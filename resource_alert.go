package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var alertSchema = map[string]*schema.Schema{
	"virtualappliance": attribute(required, vapp, forceNew),
	"name":             attribute(required, text),
	"description":      attribute(required, text),
	"alarms":           attribute(required, set(attribute(href), schema.HashString), min(1)),
}

func alertNew(d *resourceData) core.Resource {
	alarms := make([]*core.Link, 0)
	for _, a := range d.set("alarms").List() {
		alarms = append(alarms, core.NewLinkType(a.(string), "alarm"))
	}

	return &abiquo.Alert{
		Name:        d.string("name"),
		Description: d.string("description"),
		Alarms:      alarms,
	}
}

func alertEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("virtualappliance")+"/alerts", "alert")
}

func alertRead(d *resourceData, resource core.Resource) (err error) {
	alert := resource.(*abiquo.Alert)
	alarms := schema.NewSet(schema.HashString, nil)
	for _, a := range alert.Alarms {
		alarms.Add(a.URL())
	}
	d.Set("alarms", alarms)
	d.Set("name", alert.Name)
	d.Set("description", alert.Description)
	return
}
