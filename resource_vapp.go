package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

func vapp(s *schema.Schema) {
	link(s, []string{"/cloud/virtualdatacenters/[0-9]+/virtualappliances/[0-9]+$"})
}

var vappSchema = map[string]*schema.Schema{
	"name":              attribute(required, text),
	"virtualdatacenter": attribute(required, vdc, forceNew),
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
