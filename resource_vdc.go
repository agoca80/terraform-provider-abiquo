package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

func vdc(s *schema.Schema) {
	link(s, []string{"/cloud/virtualdatacenters/[0-9]+$"})
}

func vdcIP(s *schema.Schema) {
	link(s, []string{
		"/admin/enterprises/[0-9]+/limits/[0-9]+/externalnetworks/[0-9]+/ips/[0-9]+",
		"/cloud/virtualdatacenters/[0-9]+/privatenetworks/[0-9]+/ips/[0-9]+$",
		"/cloud/virtualdatacenters/[0-9]+/publicips/purchased/[0-9]+",
	})
}

var vdcSchema = map[string]*schema.Schema{
	"cpuhard":     attribute(optional, natural),
	"cpusoft":     attribute(optional, natural),
	"diskhard":    attribute(optional, natural),
	"disksoft":    attribute(optional, natural),
	"name":        attribute(required, text),
	"publichard":  attribute(optional, natural),
	"publicsoft":  attribute(optional, natural),
	"ramhard":     attribute(optional, natural),
	"ramsoft":     attribute(optional, natural),
	"storagehard": attribute(optional, natural),
	"storagesoft": attribute(optional, natural),
	"vlanhard":    attribute(optional, natural),
	"vlansoft":    attribute(optional, natural),
	"volsoft":     attribute(optional, natural),
	"volhard":     attribute(optional, natural),
	"type":        attribute(required, label([]string{"VMX_04", "KVM"}), forceNew),
	// Links
	"enterprise": attribute(required, enterprise, forceNew),
	"location":   attribute(required, location, forceNew),
	"publicips":  attribute(optional, set(attribute(ip), schema.HashString)),
	// Computed links
	"externalips":      attribute(computed, text),
	"externalnetworks": attribute(computed, text),
	"privatenetworks":  attribute(computed, text),
	"topurchase":       attribute(computed, text),
	"purchased":        attribute(computed, text),
}

func purchaseIPs(vdc core.Resource, ips *schema.Set) (err error) {
	if ips == nil {
		return
	}

	available := vdc.Rel("topurchase").Collection(nil).List()
	for _, a := range available {
		if ips.Contains(a.(*abiquo.IP).IP) {
			if err = core.Update(a.Rel("purchase"), nil); err != nil {
				break
			}
		}
	}
	return
}

func releaseIPs(resource core.Resource, ips *schema.Set) (err error) {
	purchased := resource.Rel("purchased").Collection(nil).List()
	for _, p := range purchased {
		if ips == nil || !ips.Contains(p.(*abiquo.IP).IP) {
			if err = core.Update(p.Rel("release"), nil); err != nil {
				break
			}
		}
	}
	return
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

func vdcCreate(d *resourceData, resource core.Resource) (err error) {
	// Computed links
	d.Set("externalips", resource.Rel("externalips").Href)
	d.Set("externalnetworks", resource.Rel("externalnetworks").Href)
	d.Set("privatenetworks", resource.Rel("privatenetworks").Href)
	d.Set("topurchase", resource.Rel("topurchase").Href)
	d.Set("purchased", resource.Rel("purchased").Href)
	purchaseIPs(resource, d.set("publicips"))
	return
}

func vdcUpdate(d *resourceData, resource core.Resource) (err error) {
	if err = purchaseIPs(resource, d.set("publicips")); err == nil {
		err = releaseIPs(resource, d.set("publicips"))
	}
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
	// publicips
	publicips := schema.NewSet(schema.HashString, nil)
	purchased := vdc.Rel("purchased").Collection(nil).List()
	for _, resource := range purchased {
		debug.Println("vdcRead: purchased ", resource.(*abiquo.IP).IP)
		publicips.Add(resource.(*abiquo.IP).IP)
	}
	d.Set("publicips", publicips)

	return
}

func vdcDevice(link *core.Link) (device core.Resource) {
	if vdc := link.Walk(); vdc != nil {
		device = vdc.Walk("device")
	}
	return
}
