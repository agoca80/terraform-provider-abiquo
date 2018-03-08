package main

import (
	"time"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var sgResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name": &schema.Schema{
			Required: true,
			Type:     schema.TypeString,
		},
		"cooldown": &schema.Schema{
			Required: true,
			Type:     schema.TypeInt,
		},
		"min": &schema.Schema{
			Required: true,
			Type:     schema.TypeInt,
		},
		"max": &schema.Schema{
			Required: true,
			Type:     schema.TypeInt,
		},
		"scale_out": &schema.Schema{
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"starttime": &schema.Schema{
						Optional:     true,
						Type:         schema.TypeInt,
						ValidateFunc: validateTS,
					},
					"stoptime": &schema.Schema{
						Optional:     true,
						Type:         schema.TypeInt,
						ValidateFunc: validateTS,
					},
					"numberofinstances": &schema.Schema{
						Required: true,
						Type:     schema.TypeInt,
					},
				},
			},
			MinItems: 1,
			Required: true,
			Type:     schema.TypeList,
		},
		"scale_in": &schema.Schema{
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"starttime": &schema.Schema{
						Optional:     true,
						Type:         schema.TypeInt,
						ValidateFunc: validateTS,
					},
					"stoptime": &schema.Schema{
						Optional:     true,
						Type:         schema.TypeInt,
						ValidateFunc: validateTS,
					},
					"numberofinstances": &schema.Schema{
						Required: true,
						Type:     schema.TypeInt,
					},
				},
			},
			MinItems: 1,
			Required: true,
			Type:     schema.TypeList,
		},
		"virtualappliance": &schema.Schema{
			ForceNew:     true,
			Required:     true,
			Type:         schema.TypeString,
			ValidateFunc: validateURL,
		},
		"mastervirtualmachine": &schema.Schema{
			ForceNew:     true,
			Required:     true,
			Type:         schema.TypeString,
			ValidateFunc: validateURL,
		},
	},
	Delete: sgDelete,
	Update: resourceUpdate(sgNew, "scalinggroup"),
	Create: resourceCreate(sgNew, nil, sgRead, sgEndpoint),
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
