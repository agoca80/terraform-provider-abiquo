package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var vappSchema = map[string]*schema.Schema{
	"virtualdatacenter": endpoint("virtualdatacenter"),
	"name":              attribute(required, text),
}

func vappNew(d *resourceData) core.Resource {
	return &abiquo.VirtualAppliance{
		Name: d.string("name"),
		DTO: core.NewDTO(
			d.link("virtualdatacenter"),
		),
	}
}

func vappRead(d *resourceData, resource core.Resource) (err error) {
	vapp := resource.(*abiquo.VirtualAppliance)
	d.Set("name", vapp.Name)
	return
}

var virtualappliance = &description{
	Resource: &schema.Resource{Schema: vappSchema},
	dto:      vappNew,
	endpoint: endpointPath("virtualdatacenter", "/virtualappliances"),
	media:    "virtualappliance",
	read:     vappRead,
}
