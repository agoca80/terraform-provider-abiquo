package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var actionPlanSchema = map[string]*schema.Schema{
	"virtualmachine": endpoint("virtualmachine"),
	"name":           attribute(required, text),
	"description":    attribute(required, text),
	"triggers":       attribute(optional, list(link("alarm"))),
	"entries": attribute(required, min(1), list(&schema.Resource{
		Schema: map[string]*schema.Schema{
			"parameter":     attribute(optional, text),
			"parametertype": attribute(optional, text),
			"type":          attribute(required, text),
		},
	})),
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

func actionPlanCreate(d *resourceData, resource core.Resource) (err error) {
	a := resource.(*abiquo.ActionPlan)
	if d.HasChange("triggers") {
		triggers := new(core.DTO)
		for _, trigger := range d.slice("triggers") {
			link := core.NewLinkType(trigger.(string), "alert").SetRel("alert")
			triggers.Links = append(triggers.Links, link)
		}
		err = core.Create(a.Rel("alerttriggers"), triggers)
	}
	return
}

var virtualmachineactionplan = &description{
	media:    "virtualmachineactionplan",
	dto:      actionPlanNew,
	read:     actionPlanRead,
	endpoint: endpointPath("virtualmachine", "/actionplans"),
	Resource: &schema.Resource{Schema: actionPlanSchema},
}
