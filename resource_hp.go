package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var hpResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"active":     Required().Bool(),
		"name":       Required().String(),
		"cpu":        Required().Number(),
		"ram":        Required().Number(),
		"datacenter": Required().Renew().Link(),
	},
	Delete: resourceDelete,
	Exists: resourceExists("hardwareprofile"),
	Create: resourceCreate(hpNew, nil, hpRead, hpEndpoint),
	Update: resourceUpdate(hpNew, nil, "hardwareprofile"),
	Read:   resourceRead(hpNew, hpRead, "hardwareprofile"),
}

func hpNew(d *resourceData) core.Resource {
	return &abiquo.HardwareProfile{
		Active:  d.bool("active"),
		Name:    d.string("name"),
		CPU:     d.int("cpu"),
		RAMInMB: d.int("ram"),
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
