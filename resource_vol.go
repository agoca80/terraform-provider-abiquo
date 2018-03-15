package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

var volumeSchema = map[string]*schema.Schema{
	"size": &schema.Schema{
		Required: true,
		Type:     schema.TypeInt,
	},
	"name": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"bootable": &schema.Schema{
		Optional: true,
		Type:     schema.TypeBool,
	},
	"description": &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	},
	"ctrl": &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	},
	"type": &schema.Schema{
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validation.StringInSlice([]string{"IDE", "SCSI", "VIRTIO"}, false),
	},
	"tier": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
	"virtualdatacenter": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
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
