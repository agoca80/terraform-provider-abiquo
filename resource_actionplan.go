package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var actionPlanSchema = map[string]*schema.Schema{
	"virtualmachine": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
	"name": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"description": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"triggers": &schema.Schema{
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validateURL,
		},
		Optional: true,
		Type:     schema.TypeList,
	},
	"entries": &schema.Schema{
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"parameter": &schema.Schema{
					Optional: true,
					Type:     schema.TypeString,
				},
				"parametertype": &schema.Schema{
					Optional: true,
					Type:     schema.TypeString,
				},
				"type": &schema.Schema{
					Required: true,
					Type:     schema.TypeString,
				},
			},
		},
		MinItems: 1,
		Required: true,
		Type:     schema.TypeList,
	},
}

func actionPlanNew(d *resourceData) core.Resource {
	slice := d.slice("entries")
	entries := make([]abiquo.ActionPlanEntry, len(slice))
	for i, entry := range slice {
		mapDecoder(entry, &entries[i])
		entries[i].Sequence = i
	}

	return &abiquo.ActionPlan{
		Name:        d.string("name"),
		Description: d.string("description"),
		Entries:     entries,
	}
}

func actionPlanEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("virtualmachine")+"/actionplans", "virtualmachineactionplan")
}

func actionPlanRead(d *resourceData, resource core.Resource) (err error) {
	actionPlan := resource.(*abiquo.ActionPlan)
	entries := make([]map[string]interface{}, len(actionPlan.Entries))
	for i, e := range actionPlan.Entries {
		entries[i] = map[string]interface{}{
			"parameter":     e.Parameter,
			"parametertype": e.ParameterType,
			"type":          e.Type,
		}
	}
	d.Set("name", actionPlan.Name)
	d.Set("description", actionPlan.Description)
	d.Set("entries", entries)
	return
}

func actionPlanTriggers(d *resourceData) (triggers core.DTO) {
	for _, trigger := range d.slice("triggers") {
		link := core.NewLinkType(trigger.(string), "alert").SetRel("alert")
		triggers.Links = append(triggers.Links, link)
	}
	return
}

func actionPlanCreate(d *resourceData, resource core.Resource) (err error) {
	a := resource.(*abiquo.ActionPlan)
	if triggers := actionPlanTriggers(d); len(triggers.Links) > 0 {
		err = a.SetTriggers(&triggers)
	}
	return
}
