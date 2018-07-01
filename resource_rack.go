package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var rackSchema = map[string]*schema.Schema{
	"name":        attribute(required, text),
	"number":      attribute(computed, integer), // ABICLOUDPREMIUM-10197
	"description": attribute(optional, text),
	"vlanmax":     attribute(optional, natural),
	"vlanmin":     attribute(optional, natural),
	"datacenter":  attribute(required, datacenter, forceNew),
}

func rackNew(d *resourceData) core.Resource {
	rack := &abiquo.Rack{
		ID:   d.integer("number"),
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
