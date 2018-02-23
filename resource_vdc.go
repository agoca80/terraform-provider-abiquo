package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var vdcResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name": Required().String(),
		"type": Required().Renew().String(),
		// Soft limits
		"cpusoft":     Optional().Number(),
		"disksoft":    Optional().Number(),
		"publicsoft":  Optional().Number(),
		"ramsoft":     Optional().Number(),
		"storagesoft": Optional().Number(),
		"volsoft":     Optional().Number(),
		"vlansoft":    Optional().Number(),
		// Hard limits
		"cpuhard":     Optional().Number(),
		"diskhard":    Optional().Number(),
		"publichard":  Optional().Number(),
		"ramhard":     Optional().Number(),
		"storagehard": Optional().Number(),
		"vlanhard":    Optional().Number(),
		"volhard":     Optional().Number(),
		// Links
		"enterprise": Required().Renew().Link(),
		"location":   Required().Renew().Link(),
	},
	Delete: resourceDelete,
	Create: resourceCreate(vdcNew, nil, vdcRead, vdcEndpoint),
	Exists: resourceExists("virtualdatacenter"),
	Update: resourceUpdate(vdcNew, "virtualdatacenter"),
	Read:   resourceRead(vdcNew, vdcRead, "virtualdatacenter"),
}

func vdcNew(d *resourceData) core.Resource {
	return &abiquo.VirtualDatacenter{
		Name:   d.string("name"),
		HVType: d.string("type"),
		Network: &abiquo.Network{
			Mask:    24,
			Address: "192.168.0.0",
			Gateway: "192.168.0.1",
			Name:    "tf default network",
			TypeNet: "INTERNAL",
		},
		// Soft limits
		CPUSoft:     d.int("cpusoft"),
		DiskSoft:    d.int("disksoft"),
		PublicSoft:  d.int("publicsoft"),
		RAMSoft:     d.int("ramsoft"),
		StorageSoft: d.int("storagesoft"),
		// Hard limits
		CPUHard:     d.int("cpuhard"),
		DiskHard:    d.int("diskhard"),
		PublicHard:  d.int("iphard"),
		RAMHard:     d.int("ramhard"),
		StorageHard: d.int("storagehard"),
		VLANHard:    d.int("vlanhard"),
		VLANSoft:    d.int("vlansoft"),
		DTO: core.NewDTO(
			d.link("enterprise"),
			d.link("location"),
		),
	}
}

func vdcEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType("cloud/virtualdatacenters", "virtualdatacenter")
}

func vdcRead(d *resourceData, resource core.Resource) (err error) {
	vdc := resource.(*abiquo.VirtualDatacenter)
	d.Set("name", vdc.Name)
	// Soft limits
	d.Set("cpusoft", vdc.CPUSoft)
	d.Set("disksoft", vdc.DiskSoft)
	d.Set("publicsoft", vdc.PublicHard)
	d.Set("ramsoft", vdc.RAMSoft)
	d.Set("storagesoft", vdc.StorageHard)
	d.Set("vlansoft", vdc.VLANSoft)
	// Hard limits
	d.Set("cpuhard", vdc.CPUHard)
	d.Set("diskhard", vdc.DiskHard)
	d.Set("publichard", vdc.PublicHard)
	d.Set("ramhard", vdc.RAMSoft)
	d.Set("storagehard", vdc.StorageHard)
	d.Set("vlanhard", vdc.VLANSoft)
	return
}
