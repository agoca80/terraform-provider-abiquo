package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var deviceSchema = map[string]*schema.Schema{
	"datacenter":  endpoint("datacenter"),
	"description": attribute(optional, text),
	"endpoint":    attribute(required, forceNew, href),
	"name":        attribute(required, text),
	"password":    attribute(optional, text, sensitive),
	"username":    attribute(optional, text),
	"default":     attribute(optional, boolean),
	"devicetype":  attribute(required, forceNew, link("devicetype")),
	"enterprise":  attribute(optional, forceNew, link("enterprise")),
}

func deviceDTO(d *resourceData) core.Resource {
	return &abiquo.Device{
		Description: d.string("description"),
		Endpoint:    d.string("endpoint"),
		Name:        d.string("name"),
		Username:    d.string("username"),
		Password:    d.string("password"),
		Default:     d.boolean("default"),
		DTO: core.NewDTO(
			d.link("enterprise"),
			d.link("devicetype"),
		),
	}
}

func deviceRead(d *resourceData, resource core.Resource) (err error) {
	device := resource.(*abiquo.Device)
	d.Set("endpoint", device.Endpoint)
	d.Set("name", device.Name)
	d.SetOk("password", device.Password)
	d.SetOk("username", device.Username)
	d.SetOk("description", device.Description)
	d.SetOk("enterprise", device.Rel("enterprise").URL())
	return
}

var device = &description{
	Resource: &schema.Resource{Schema: deviceSchema},
	dto:      deviceDTO,
	endpoint: endpointPath("datacenter", "/devices"),
	media:    "device",
	read:     deviceRead,
}
