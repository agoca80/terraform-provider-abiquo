package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var hardwareprofileSchema = map[string]*schema.Schema{
	"datacenter": endpoint("datacenter"),
	"active":     attribute(required, boolean),
	"name":       attribute(required, text),
	"cpu":        attribute(required, natural),
	"ram":        attribute(required, natural),
}

func hardwareprofileNew(d *resourceData) core.Resource {
	return &abiquo.HardwareProfile{
		Active:  d.boolean("active"),
		Name:    d.string("name"),
		CPU:     d.integer("cpu"),
		RAMInMB: d.integer("ram"),
	}
}

func hardwareprofileRead(d *resourceData, resource core.Resource) (err error) {
	hardwareprofile := resource.(*abiquo.HardwareProfile)
	d.Set("active", hardwareprofile.Active)
	d.Set("name", hardwareprofile.Name)
	d.Set("cpu", hardwareprofile.CPU)
	d.Set("ram", hardwareprofile.RAMInMB)
	d.Set("datacenter", hardwareprofile.Rel("datacenter").URL())
	return
}

var hardwareprofile = &description{
	media:    "hardwareprofile",
	dto:      hardwareprofileNew,
	read:     hardwareprofileRead,
	endpoint: endpointPath("datacenter", "/hardwareprofiles"),
	Resource: &schema.Resource{Schema: hardwareprofileSchema},
}
