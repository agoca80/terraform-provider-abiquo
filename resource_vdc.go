package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

var vdcSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"type": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validation.StringInSlice([]string{"VMX_04", "KVM"}, false),
	},
	// Soft limits
	"cpusoft": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	"disksoft": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	"publicsoft": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	"ramsoft": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	"storagesoft": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	"volsoft": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	"vlansoft": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	// Hard limits
	"cpuhard": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	"diskhard": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	"publichard": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	"ramhard": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	"storagehard": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	"vlanhard": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	"volhard": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	// Links
	"enterprise": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
	"location": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
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

func vdcDevice(link *core.Link) (device core.Resource) {
	if vdc := link.Walk(); vdc != nil {
		device = vdc.Walk("device")
	}
	return
}
