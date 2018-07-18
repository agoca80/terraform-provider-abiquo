package main

import (
	"fmt"
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

func sgEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("virtualappliance")+"/scalinggroups", "scalinggroup")
}

func sgRead(d *resourceData, resource core.Resource) (e error) {
	sg := resource.(*abiquo.ScalingGroup)
	d.Set("name", sg.Name)
	d.Set("cooldown", sg.Cooldown)
	d.Set("max", sg.Max)
	d.Set("min", sg.Min)
	return
}

func sgUpdate(rd *schema.ResourceData, m interface{}) (err error) {
	d := newResourceData(rd, "scalinggroup")
	resource := d.Link.Walk()
	if resource == nil {
		return fmt.Errorf("scaling group %q was not found", d.Id())
	}

	sg := resource.(*abiquo.ScalingGroup)
	if !sg.Maintenance {
		if err = sg.StartMaintenance(); err != nil {
			return
		}
	}

	// Update the SG
	modify := sgNew(d).(*abiquo.ScalingGroup)
	if err = core.Update(d, modify); err != nil {
		return
	}

	err = modify.EndMaintenance()

	return
}

func sgDelete(rd *schema.ResourceData, m interface{}) (err error) {
	d := newResourceData(rd, "scalinggroup")
	resource := d.Link.Walk()
	if resource == nil {
		return fmt.Errorf("scaling group %q was not found", d.Id())
	}

	sg := resource.(*abiquo.ScalingGroup)
	if !sg.Maintenance {
		if err = sg.StartMaintenance(); err != nil {
			return
		}
	}

	// Delete the SG
	if err = core.Remove(sg); err != nil {
		return
	}

	// Delete the SG VMs
	vms := sg.Links.Filter(func(l *core.Link) bool {
		return l.Rel == "virtualmachine"
	})
	for _, vm := range vms {
		if err = core.Remove(vm); err != nil {
			return
		}
	}

	return
}
