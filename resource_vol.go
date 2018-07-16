package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var volumeSchema = map[string]*schema.Schema{
	"virtualdatacenter": endpoint("virtualdatacenter"),
	"size":              attribute(required, natural),
	"name":              attribute(required, text),
	"bootable":          attribute(optional, boolean),
	"description":       attribute(optional, text),
	"ctrl":              attribute(optional, text),
	"type":              attribute(required, label([]string{"IDE", "SCSI", "VIRTIO"})),
	"tier":              attribute(required, forceNew, href),
}

func volNew(d *resourceData) core.Resource {
	return &abiquo.Volume{
		Name:               d.string("name"),
		Description:        d.string("description"),
		DiskControllerType: d.string("type"),
		DiskController:     d.string("ctrl"),
		Bootable:           d.boolean("bootable"),
		SizeInMB:           d.integer("size"),
		DTO: core.NewDTO(
			d.link("tier"),
		),
	}
}

func volRead(d *resourceData, resource core.Resource) (e error) {
	v := resource.(*abiquo.Volume)
	d.Set("name", v.Name)
	d.Set("bootable", v.Bootable)
	d.Set("description", v.Description)
	d.Set("type", v.DiskControllerType)
	d.Set("ctrl", v.DiskController)
	d.Set("size", v.SizeInMB)
	d.Set("tier", v.Rel("tier").URL())
	return
}

func volEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("virtualdatacenter")+"/volumes", "volume")
}
