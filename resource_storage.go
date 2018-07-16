package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var storageDeviceSchema = map[string]*schema.Schema{
	"datacenter": endpoint("datacenter"),
	"name":       attribute(required, text),
	"ip":         attribute(required, ip),
}

func storageDeviceNew(d *resourceData) core.Resource {
	return &abiquo.StorageDevice{
		Name:           d.string("name"),
		Technology:     "NFS",
		ManagementIP:   d.string("ip"),
		ManagementPort: 2049,
		ServiceIP:      d.string("ip"),
		ServicePort:    2049,
		DTO: core.NewDTO(
			d.link("datacenter"),
		),
	}
}

func storageDeviceEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("datacenter")+"/storage/devices", "storagedevice")
}

func storageDeviceRead(d *resourceData, resource core.Resource) (err error) {
	storageDevice := resource.(*abiquo.StorageDevice)
	d.Set("name", storageDevice.Name)
	return
}
