package main

import (
	"time"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var sgScaleResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"starttime":         attribute(optional, timestamp),
		"stoptime":          attribute(optional, timestamp),
		"numberofinstances": attribute(required, natural),
	},
}

var sgSchema = map[string]*schema.Schema{
	"name":                 attribute(required, text),
	"cooldown":             attribute(required, natural),
	"min":                  attribute(required, natural),
	"max":                  attribute(required, natural),
	"scale_out":            attribute(required, list(sgScaleResource)),
	"scale_in":             attribute(required, list(sgScaleResource)),
	"virtualappliance":     attribute(required, forceNew, link("virtualappliance")),
	"mastervirtualmachine": attribute(required, forceNew, href),
}

func sgRules(rules []interface{}) (sgRules []abiquo.ScalingGroupRule) {
	for _, r := range rules {
		rule := struct {
			NumberOfInstances int
			StartTime         string
			EndType           string
		}{}
		mapDecoder(r, &rule)
		from, _ := time.Parse(tsFormat, rule.StartTime)
		until, _ := time.Parse(tsFormat, rule.EndType)
		sgRules = append(sgRules, abiquo.ScalingGroupRule{
			NumberOfInstances: rule.NumberOfInstances,
			StartTime:         from.Unix(),
			EndTime:           until.Unix(),
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
	return
}

func sgDelete(rd *schema.ResourceData, m interface{}) (err error) {
	d := newResourceData(rd, "scalinggroup")
	// Get the VMs inside the SG
	sg := new(abiquo.ScalingGroup)
	if err = core.Read(d, sg); err != nil {
		return
	}

	vms := []*core.Link{}
	for _, l := range sg.Links {
		if l.Rel == "virtualmachine" {
			vms = append(vms, l)
		}
	}

	// Go to maintenance mode
	if !sg.Maintenance {
		if err = sg.StartMaintenance(); err != nil {
			return
		}
	}

	// PENDING move to ojal/abiquo
	// Delete the SG
	if err = core.Remove(d); err == nil {
		// Delete the SG VMs
		for _, vm := range vms {
			if err = core.Remove(vm); err != nil {
				return
			}
		}
	}
	return
}
