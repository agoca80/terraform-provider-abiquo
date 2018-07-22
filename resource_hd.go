package main

import (
	"strings"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var hdSchema = map[string]*schema.Schema{
	"virtualdatacenter": endpoint("virtualdatacenter"),
	"size":              attribute(required, natural),
	"label":             attribute(required, text),
	"type":              attribute(required, label([]string{"IDE", "SCSI", "VIRTIO"})),
	"ctrl":              attribute(optional, text),
	"dstier":            attribute(optional, computed, href),
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
		SizeInMb:           d.integer("size"),
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

var harddisk = &description{
	dto:      hdNew,
	endpoint: endpointPath("virtualdatacenter", "/disks"),
	media:    "harddisk",
	read:     hdRead,
	Resource: &schema.Resource{
		Schema: hdSchema,
		Update: schema.Noop,
		Delete: schema.Noop,
	},
}
