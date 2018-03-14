package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var deviceSchema = map[string]*schema.Schema{
	"devicetype": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
	"endpoint": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
	"description": &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	},
	"name": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"password": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"username": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"enterprise": &schema.Schema{
		ForceNew:     true,
		Optional:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
	"datacenter": &schema.Schema{
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
}

var deviceResource = &schema.Resource{
	Schema: deviceSchema,
	Delete: resourceDelete,
	Exists: resourceExists("device"),
	Create: resourceCreate(deviceDTO, nil, deviceRead, deviceEndpoint),
	Update: resourceUpdate(deviceDTO, nil, "device"),
	Read:   resourceRead(deviceDTO, deviceRead, "device"),
}

func deviceDTO(d *resourceData) core.Resource {
	return &abiquo.Device{
		Description: d.string("description"),
		Endpoint:    d.string("endpoint"),
		Name:        d.string("name"),
		Username:    d.string("username"),
		Password:    d.string("password"),
		Default:     true,
		DTO: core.NewDTO(
			d.link("enterprise"),
			d.link("devicetype"),
		),
	}
}

func deviceEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("datacenter")+"/devices", "device")
}

func deviceRead(d *resourceData, resource core.Resource) (err error) {
	device := resource.(*abiquo.Device)
	d.Set("endpoint", device.Endpoint)
	d.Set("name", device.Name)
	d.Set("password", device.Password)
	d.Set("username", device.Username)

	if _, ok := d.GetOk("description"); ok {
		d.Set("description", device.Description)
	}

	if _, ok := d.GetOk("enterprise"); ok {
		d.Set("enterprise", device.Rel("enterprise").URL())
	}

	return
}
