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

func vappEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("virtualdatacenter")+"/virtualappliances", "virtualappliance")
}

func vappRead(d *resourceData, resource core.Resource) (err error) {
	vapp := resource.(*abiquo.VirtualAppliance)
	d.Set("name", vapp.Name)
	return
}

var resourceVapp = &schema.Resource{
	Schema: vappSchema,
	Delete: resourceDelete,
	Exists: resourceExists("virtualappliance"),
	Create: resourceCreate(vappNew, nil, vappRead, vappEndpoint),
	Update: resourceUpdate(vappNew, nil, "virtualappliance"),
	Read:   resourceRead(vappNew, vappRead, "virtualappliance"),
}
