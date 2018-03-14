package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var vappSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"virtualdatacenter": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
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
