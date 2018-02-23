package main

import (
	"strings"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var hdResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"size":              Required().Number(),
		"label":             Required().String(),
		"type":              Required().ValidateString([]string{"IDE", "SCSI", "VIRTIO"}),
		"ctrl":              Optional().String(),
		"dstier":            Optional().Renew().Link(),
		"virtualdatacenter": Required().Renew().Link(),
	},
	Delete: hdDelete,
	Create: resourceCreate(hdNew, nil, hdRead, hdEndpoint),
	Exists: resourceExists("harddisk"),
	Update: hdUpdate,
	Read:   resourceRead(hdNew, hdRead, "harddisk"),
}

func hdLink(href string) *core.Link {
	if strings.Contains(href, "/disks/") {
		return core.NewLinkType(href, "harddisk")
	}
	return core.NewLinkType(href, "volume")
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

func hdDelete(d *schema.ResourceData, m interface{}) (err error) {
	return
}

func hdUpdate(d *schema.ResourceData, m interface{}) (err error) {
	return
}

func hdEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("virtualdatacenter")+"/disks", "harddisk")
}
