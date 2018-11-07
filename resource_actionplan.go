package main

import (
	"strings"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var actionPlanEntryResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"links":         attribute(optional, list(href)),
		"parameter":     attribute(optional, text),
		"parametertype": attribute(optional, text),
		"type":          attribute(required, text),
	},
}

var actionPlanSchema = map[string]*schema.Schema{
	"name":        attribute(required, text),
	"description": attribute(required, text),
	"triggers":    attribute(optional, list(link("alarm"))),
	"entries":     attribute(required, min(1), list(actionPlanEntryResource)),
}

func actionPlanEntries(d *resourceData) (entries []abiquo.ActionPlanEntry) {
	for sequence, e := range d.slice("entries") {
		dto := core.NewDTO()
		entry := e.(map[string]interface{})
		for _, h := range entry["links"].([]interface{}) {
			href := h.(string)
			if strings.Contains(href, "/virtualmachines/") {
				dto.Add(core.NewLinkType(href, "virtualmachine").SetRel("virtualmachine"))
			} else {
				panic("Not implemented")
			}
		}
		entries = append(entries, abiquo.ActionPlanEntry{
			Sequence:      sequence,
			Parameter:     entry["parameter"].(string),
			ParameterType: entry["parametertype"].(string),
			Type:          entry["type"].(string),
			DTO:           dto,
		})
	}
	return
}

func actionPlanNew(d *resourceData) core.Resource {
	return &abiquo.ActionPlan{
		Name:        d.string("name"),
		Description: d.string("description"),
		Entries:     actionPlanEntries(d),
	}
}

func actionPlanRead(d *resourceData, resource core.Resource) (err error) {
	actionPlan := resource.(*abiquo.ActionPlan)
	entries := make([]map[string]interface{}, len(actionPlan.Entries))
	for i, e := range actionPlan.Entries {
		links := []interface{}{}
		for _, h := range e.Links {
			links = append(links, h.URL())
		}
		entries[i] = map[string]interface{}{
			"parameter":     e.Parameter,
			"parametertype": e.ParameterType,
			"type":          e.Type,
			"links":         links,
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

var actionplan = &description{
	media:    "actionplan",
	dto:      actionPlanNew,
	read:     actionPlanRead,
	endpoint: endpointConst("cloud/actionplans"),
	Resource: &schema.Resource{Schema: actionPlanSchema},
}
