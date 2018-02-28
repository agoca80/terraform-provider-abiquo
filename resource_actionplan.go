package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var actionPlanEntrySchema = map[string]*schema.Schema{
	"parameter":     Optional().String(),
	"parametertype": Optional().String(),
	"type":          Required().String(),
}

var actionPlanSchema = map[string]*schema.Schema{
	"virtualmachine": Required().Renew().Link(),
	"name":           Required().String(),
	"description":    Required().String(),
	"triggers":       Optional().Links(),
	"entries": Required().list(1, &schema.Resource{
		Schema: actionPlanEntrySchema,
	}),
}

var actionPlanResource = &schema.Resource{
	Schema: actionPlanSchema,
	Delete: resourceDelete,
	Exists: resourceExists("virtualmachineactionplan"),
	Update: resourceUpdate(actionPlanNew, "virtualmachineactionplan"),
	Create: actionPlanCreate,
	Read:   resourceRead(actionPlanNew, actionPlanRead, "virtualmachineactionplan"),
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

func actionPlanCreate(rd *schema.ResourceData, m interface{}) (err error) {
	d := newResourceData(rd, "")
	plan := actionPlanNew(d).(*abiquo.ActionPlan)

	if err = core.Create(actionPlanEndpoint(d), plan); err == nil {
		d.SetId(plan.URL())
		if triggers := actionPlanTriggers(d); len(triggers.Links) > 0 {
			err = plan.SetTriggers(&triggers)
		}
	}
	return
}
