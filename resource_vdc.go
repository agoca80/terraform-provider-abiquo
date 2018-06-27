package main

import (
	"fmt"

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
	// Public ips
	"publicips": &schema.Schema{
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validateIP,
		},
		Optional: true,
		Set:      schema.HashString,
		Type:     schema.TypeSet,
	},
	// Computed links
	"externalips": &schema.Schema{
		Computed: true,
		Type:     schema.TypeString,
	},
	"externalnetworks": &schema.Schema{
		Computed: true,
		Type:     schema.TypeString,
	},
	"privatenetworks": &schema.Schema{
		Computed: true,
		Type:     schema.TypeString,
	},
	"topurchase": &schema.Schema{
		Computed: true,
		Type:     schema.TypeString,
	},
	"purchased": &schema.Schema{
		Computed: true,
		Type:     schema.TypeString,
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
			Name:    "default",
			Type:    "INTERNAL",
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

func purchaseIP(available core.Resources, address string) (err error) {
	if ip := available.Find(func(r core.Resource) (found bool) {
		return r.(*abiquo.IP).IP == address
	}); ip != nil {
		return core.Update(ip.Rel("purchase"), nil)
	}
	return fmt.Errorf("ip %q was not found", address)
}

func vdcCreate(d *resourceData, resource core.Resource) (err error) {
	// Computed links
	d.Set("externalips", resource.Rel("externalips").Href)
	d.Set("externalnetworks", resource.Rel("externalnetworks").Href)
	d.Set("privatenetworks", resource.Rel("privatenetworks").Href)
	d.Set("topurchase", resource.Rel("topurchase").Href)
	d.Set("purchased", resource.Rel("purchased").Href)

	// Buy public IPs
	ipsSet := d.set("publicips")
	if ipsSet == nil || ipsSet.Len() == 0 {
		return
	}

	addresses := ipsSet.List()
	available := resource.Rel("topurchase").Collection(nil).List()
	for _, address := range addresses {
		if err = purchaseIP(available, address.(string)); err != nil {
			break
		}
	}

	return
}

func vdcUpdate(d *resourceData, resource core.Resource) (err error) {
	// vdc := resource.(*abiquo.VirtualDatacenter)
	if !d.HasChange("publicips") {
		return
	}

	// release IPs

	// buy missing IPs

	return
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
