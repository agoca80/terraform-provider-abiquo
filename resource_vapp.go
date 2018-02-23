package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var vappResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name":              Required().String(),
		"virtualdatacenter": Required().Renew().Link(),
	},
	Delete: resourceDelete,
	Exists: resourceExists("virtualappliance"),
	Create: resourceCreate(vappNew, nil, vappRead, vappEndpoint),
	Update: resourceUpdate(vappNew, "virtualappliance"),
	Read:   resourceRead(vappNew, vappRead, "virtualappliance"),
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
