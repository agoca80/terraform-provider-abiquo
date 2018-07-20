package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var alertSchema = map[string]*schema.Schema{
	"virtualappliance": endpoint("virtualappliance"),
	"name":             attribute(required, text),
	"description":      attribute(optional, text),
	"subscribers":      attribute(optional, set(email)),
	"alarms":           attribute(required, set(href), min(1)),
}

func alertNew(d *resourceData) core.Resource {
	alarms := core.Links{}
	for _, a := range d.set("alarms").List() {
		alarms = append(alarms, core.NewLinkType(a.(string), "alarm"))
	}

	subscribers := []string{}
	if d.set("subscribers") != nil {
		for _, s := range d.set("subscribers").List() {
			subscribers = append(subscribers, s.(string))
		}
	}

	return &abiquo.Alert{
		Name:        d.string("name"),
		Description: d.string("description"),
		Alarms:      alarms,
		Subscribers: subscribers,
	}
}

func alertEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("virtualappliance")+"/alerts", "alert")
}

func alertRead(d *resourceData, resource core.Resource) (err error) {
	alert := resource.(*abiquo.Alert)
	alarms := []interface{}{}
	alert.Alarms.Map(func(l *core.Link) {
		alarms = append(alarms, l.URL())
	})

	d.Set("subscribers", alert.Subscribers)
	d.Set("alarms", alarms)
	d.Set("name", alert.Name)
	d.Set("description", alert.Description)
	return
}

var resourceAlert = &schema.Resource{
	Schema: alertSchema,
	Delete: resourceDelete,
	Exists: resourceExists("alert"),
	Update: resourceUpdate(alertNew, nil, "alert"),
	Create: resourceCreate(alertNew, nil, alertRead, alertEndpoint),
	Read:   resourceRead(alertNew, alertRead, "alert"),
}
