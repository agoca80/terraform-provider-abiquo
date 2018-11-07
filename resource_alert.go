package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var alertSchema = map[string]*schema.Schema{
	"name":        attribute(required, text),
	"description": attribute(optional, text),
	"subscribers": attribute(optional, set(email)),
	"alarms":      attribute(required, set(href), min(1)),
}

func alertNew(d *resourceData) core.Resource {
	alarms := core.NewDTO()
	for _, a := range d.set("alarms").List() {
		alarms.Add(core.NewLinkType(a.(string), "alarm").SetRel("alarm"))
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
		DTO:         alarms,
		Subscribers: subscribers,
	}
}

func alertRead(d *resourceData, resource core.Resource) (err error) {
	alert := resource.(*abiquo.Alert)
	alarms := []interface{}{}
	alert.Map(func(l *core.Link) {
		if l.Rel == "alarm" {
			alarms = append(alarms, l.URL())
		}
	})

	d.Set("subscribers", alert.Subscribers)
	d.Set("alarms", alarms)
	d.Set("name", alert.Name)
	d.Set("description", alert.Description)
	return
}

var alert = &description{
	Resource: &schema.Resource{Schema: alertSchema},
	dto:      alertNew,
	endpoint: endpointConst("cloud/alerts"),
	media:    "alert",
	read:     alertRead,
}
