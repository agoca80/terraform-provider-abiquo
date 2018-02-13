package main

import (
	"github.com/abiquo/opal/abiquo"
	"github.com/abiquo/opal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var volResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"size":              Required().Number(),
		"name":              Required().String(),
		"bootable":          Optional().Bool(),
		"description":       Optional().String(),
		"ctrl":              Optional().String(),
		"type":              Required().ValidateString([]string{"IDE", "SCSI", "VIRTIO"}),
		"tier":              Required().Link(),
		"virtualdatacenter": Required().Link(),
	},
	Delete: resourceDelete,
	Update: resourceUpdate(volNew, "volume"),
	Create: resourceCreate(volNew, volEndpoint),
	Read:   resourceRead(volNew, volRead, "volume"),
}

func volNew(d *resourceData) core.Resource {
	return &abiquo.Volume{
		Name:               d.string("name"),
		Description:        d.string("description"),
		DiskControllerType: d.string("type"),
		DiskController:     d.string("ctrl"),
		Bootable:           d.bool("bootable"),
		SizeInMB:           d.int("size"),
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
