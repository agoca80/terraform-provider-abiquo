package main

import (
	"time"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var sgScaleResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"endtime":           attribute(optional, timestamp, byDefault("")),
		"starttime":         attribute(optional, timestamp, byDefault("")),
		"numberofinstances": attribute(required, natural),
	},
}

var sgSchema = map[string]*schema.Schema{
	"virtualappliance":     endpoint("virtualappliance"),
	"name":                 attribute(required, text),
	"cooldown":             attribute(required, natural),
	"min":                  attribute(required, natural),
	"max":                  attribute(required, natural),
	"scale_out":            attribute(required, list(sgScaleResource)),
	"scale_in":             attribute(required, list(sgScaleResource)),
	"mastervirtualmachine": attribute(required, link("virtualmachine"), forceNew),
}

func ruleTS(dateStr string) (timestamp int64) {
	if dateStr != "" {
		date, _ := time.Parse(time.RFC3339, dateStr)
		timestamp = date.Unix()
	}
	return
}

func sgRules(rules []interface{}) (sgRules []abiquo.ScalingGroupRule) {
	for _, r := range rules {
		rule := r.(map[string]interface{})
		sgRules = append(sgRules, abiquo.ScalingGroupRule{
			NumberOfInstances: rule["numberofinstances"].(int),
			StartTime:         ruleTS(rule["starttime"].(string)),
			EndTime:           ruleTS(rule["endtime"].(string)),
		})
	}
	return
}

func sgNew(d *resourceData) core.Resource {
	return &abiquo.ScalingGroup{
		Name:     d.string("name"),
		Cooldown: d.integer("cooldown"),
		Max:      d.integer("max"),
		Min:      d.integer("min"),
		ScaleIn:  sgRules(d.slice("scale_in")),
		ScaleOut: sgRules(d.slice("scale_out")),
		DTO: core.NewDTO(
			d.link("mastervirtualmachine"),
		),
	}
}

func sgRead(d *resourceData, resource core.Resource) (e error) {
	sg := resource.(*abiquo.ScalingGroup)
	d.Set("name", sg.Name)
	d.Set("cooldown", sg.Cooldown)
	d.Set("max", sg.Max)
	d.Set("min", sg.Min)
	return
}

func sgUpdate(rd *schema.ResourceData, _ interface{}) (err error) {
	d := newDataType(rd, "scalinggroup")
	resource, err := d.Link.Walk()
	if err != nil {
		return
	}

	sg := resource.(*abiquo.ScalingGroup)
	if !sg.Maintenance {
		if err = sg.StartMaintenance(); err != nil {
			return
		}
	}

	// Update the SG
	modify := sgNew(d).(*abiquo.ScalingGroup)
	err = d.Update(modify)
	if err != nil {
		return
	}

	return modify.EndMaintenance()
}

func sgDelete(rd *schema.ResourceData, m interface{}) (err error) {
	d := newDataType(rd, "scalinggroup")
	resource, err := d.Link.Walk()
	if err != nil {
		return
	}

	sg := resource.(*abiquo.ScalingGroup)
	if !sg.Maintenance {
		if err = sg.StartMaintenance(); err != nil {
			return
		}
	}

	// Delete the SG
	err = sg.Remove()
	if err != nil {
		return
	}

	// Delete the SG VMs
	vms := sg.Links.Filter(func(l *core.Link) bool {
		return l.Rel == "virtualmachine"
	})
	for _, vm := range vms {
		err = vm.Remove()
		if err != nil {
			return
		}
	}

	return
}

var scalinggroup = &description{
	dto:      sgNew,
	endpoint: endpointPath("virtualappliance", "/scalinggroups"),
	media:    "scalinggroup",
	read:     sgRead,
	Resource: &schema.Resource{
		Schema: sgSchema,
		Delete: sgDelete,
		Update: sgUpdate,
	},
}
