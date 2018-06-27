package main

import (
	"strings"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

var hdSchema = map[string]*schema.Schema{
	"size": &schema.Schema{
		Required: true,
		Type:     schema.TypeInt,
	},
	"label": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"type": &schema.Schema{
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validation.StringInSlice([]string{"IDE", "SCSI", "VIRTIO"}, false),
	},
	"ctrl": &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	},
	"dstier": &schema.Schema{
		ForceNew:     true,
		Optional:     true,
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

func hdLink(href string) *core.Link {
	var media string
	if harddisk := strings.Contains(href, "/disks/"); harddisk {
		media = "harddisk"
	} else {
		media = "volume"
	}
	return core.NewLinkType(href, media)
}

func hdNew(d *resourceData) core.Resource {
	return &abiquo.HardDisk{
		Label:              d.string("label"),
		DiskController:     d.string("ctrl"),
		DiskControllerType: d.string("type"),
		SizeInMb:           d.int("size"),
		DTO:                core.NewDTO(d.link("dstier")),
	}
}

func hdRead(d *resourceData, resource core.Resource) (err error) {
	hd := resource.(*abiquo.HardDisk)
	d.Set("label", hd.Label)
	d.Set("type", hd.DiskControllerType)
	d.Set("size", hd.SizeInMb)
	if _, ok := d.GetOk("ctrl"); ok {
		d.Set("ctrl", hd.DiskController)
	}
	return
}

func hdEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("virtualdatacenter")+"/disks", "harddisk")
}
