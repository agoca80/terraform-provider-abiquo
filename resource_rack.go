package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var rackSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"number": &schema.Schema{
		Computed: true,
		Type:     schema.TypeInt,
	}, // ABICLOUDPREMIUM-10197
	"description": &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	},
	"vlanmax": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	"vlanmin": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	"datacenter": &schema.Schema{
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
}

func rackNew(d *resourceData) core.Resource {
	rack := &abiquo.Rack{
		ID:   d.int("number"),
		Name: d.string("name"),
	}

	if d, ok := d.GetOk("description"); ok {
		rack.Description = d.(string)
	}

	if min, ok := d.GetOk("vlanmin"); ok {
		rack.VlanIDMin = min.(int)
	}

	if max, ok := d.GetOk("vlanmax"); ok {
		rack.VlanIDMax = max.(int)
	}

	return rack
}

func rackEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("datacenter")+"/racks", "rack")
}

func rackRead(d *resourceData, resource core.Resource) (err error) {
	rack := resource.(*abiquo.Rack)

	d.Set("name", rack.Name)
	d.Set("number", rack.ID)

	if _, ok := d.GetOk("description"); ok {
		d.Set("description", rack.Description)
	}

	if _, ok := d.GetOk("vlanmin"); ok {
		d.Set("vlanmin", rack.VlanIDMin)
	}

	if _, ok := d.GetOk("vlanmax"); ok {
		d.Set("vlanmax", rack.VlanIDMax)
	}

	return
}
