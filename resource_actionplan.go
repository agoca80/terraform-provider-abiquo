package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var actionPlanEntryResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"parameter":     attribute(optional, text),
		"parametertype": attribute(optional, text),
		"type":          attribute(required, text),
	},
}

var actionPlanSchema = map[string]*schema.Schema{
	"virtualmachine": endpoint("virtualmachine"),
	"name":           attribute(required, text),
	"description":    attribute(required, text),
	"triggers":       attribute(optional, list(link("alarm"))),
	"entries":        attribute(required, min(1), list(actionPlanEntryResource)),
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

var resourcePlan = &schema.Resource{
	Schema: actionPlanSchema,
	Delete: resourceDelete,
	Exists: resourceExists("virtualmachineactionplan"),
	Update: resourceUpdate(actionPlanNew, nil, "virtualmachineactionplan"),
	Create: resourceCreate(actionPlanNew, actionPlanCreate, actionPlanRead, actionPlanEndpoint),
	Read:   resourceRead(actionPlanNew, actionPlanRead, "virtualmachineactionplan"),
}
