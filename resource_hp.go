package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var hpSchema = map[string]*schema.Schema{
	"active":     attribute(required, boolean),
	"name":       attribute(required, text),
	"cpu":        attribute(required, natural),
	"ram":        attribute(required, natural),
	"datacenter": attribute(required, datacenter, forceNew),
}

func hpNew(d *resourceData) core.Resource {
	return &abiquo.HardwareProfile{
		Active:  d.boolean("active"),
		Name:    d.string("name"),
		CPU:     d.integer("cpu"),
		RAMInMB: d.integer("ram"),
	}
}

func hpEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("datacenter")+"/hardwareprofiles", "hardwareprofile")
}

func hpRead(d *resourceData, resource core.Resource) (err error) {
	hp := resource.(*abiquo.HardwareProfile)
	d.Set("active", hp.Active)
	d.Set("name", hp.Name)
	d.Set("cpu", hp.CPU)
	d.Set("ram", hp.RAMInMB)
	d.Set("datacenter", hp.Rel("datacenter").URL())
	return
}
