package main

import (
	"time"

	"github.com/abiquo/opal/abiquo"
	"github.com/abiquo/opal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var sgResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name":     Required().String(),
		"cooldown": Required().Number(),
		"min":      Required().Number(),
		"max":      Required().Number(),
		"scale_out": required(&schema.Schema{Type: schema.TypeList, MinItems: 1, Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"starttime":         Optional().Timestamp(),
				"stoptime":          Optional().Timestamp(),
				"numberofinstances": Required().Number(),
			},
		}}),
		"scale_in": required(&schema.Schema{Type: schema.TypeList, MinItems: 1, Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"starttime":         Optional().Timestamp(),
				"stoptime":          Optional().Timestamp(),
				"numberofinstances": Required().Number(),
			},
		}}),
		"virtualappliance":     Required().Renew().Link(),
		"mastervirtualmachine": Required().Renew().Link(),
	},
	Delete: sgDelete,
	Update: resourceUpdate(sgNew, "scalinggroup"),
	Create: resourceCreate(sgNew, sgEndpoint),
	Read:   resourceRead(sgNew, sgRead, "scalinggroup"),
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
		Cooldown: d.int("cooldown"),
		Max:      d.int("max"),
		Min:      d.int("min"),
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

	// PENDING move to opal/abiquo
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
