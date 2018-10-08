package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var vdcSchema = map[string]*schema.Schema{
	"cpuhard":     attribute(optional, natural),
	"cpusoft":     attribute(optional, natural),
	"diskhard":    attribute(optional, natural),
	"disksoft":    attribute(optional, natural),
	"name":        attribute(required, text),
	"net_address": attribute(optional, ip, byDefault("172.16.0.0")),
	"net_gateway": attribute(optional, ip, byDefault("172.16.0.1")),
	"net_dns1":    attribute(optional, ip),
	"net_dns2":    attribute(optional, ip),
	"net_name":    attribute(optional, text, byDefault("private_network")),
	"net_mask":    attribute(optional, positive, byDefault(16)),
	"net_suffix":  attribute(optional, text),
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
	"type":        attribute(required, label(machineType), forceNew),
	// Links
	"enterprise": attribute(required, forceNew, link("enterprise")),
	"location":   attribute(required, forceNew, link("location")),
	"publicips":  attribute(optional, set(ip)),
	// Computed links
	"device":           attribute(computed, text),
	"externalips":      attribute(computed, text),
	"externalnetworks": attribute(computed, text),
	"network":          attribute(computed, text),
	"privatenetworks":  attribute(computed, text),
	"templates":        attribute(computed, text),
	"topurchase":       attribute(computed, text),
	"purchased":        attribute(computed, text),
	"tiers":            attribute(computed, text),
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
			Address: d.string("net_address"),
			DNS1:    d.string("net_dns1"),
			DNS2:    d.string("net_dns2"),
			Gateway: d.string("net_gateway"),
			Mask:    d.integer("net_mask"),
			Name:    d.string("net_name"),
			Suffix:  d.string("net_suffix"),
			Type:    "INTERNAL",
		},
		// Soft limits
		CPUSoft:     d.integer("cpusoft"),
		DiskSoft:    d.integer("disksoft"),
		PublicSoft:  d.integer("publicsoft"),
		RAMSoft:     d.integer("ramsoft"),
		StorageSoft: d.integer("storagesoft"),
		// Hard limits
		CPUHard:     d.integer("cpuhard"),
		DiskHard:    d.integer("diskhard"),
		PublicHard:  d.integer("iphard"),
		RAMHard:     d.integer("ramhard"),
		StorageHard: d.integer("storagehard"),
		VLANHard:    d.integer("vlanhard"),
		VLANSoft:    d.integer("vlansoft"),
		DTO: core.NewDTO(
			d.link("enterprise"),
			d.link("location"),
		),
	}
}

func vdcCreate(d *resourceData, resource core.Resource) (err error) {
	d.Set("device", resource.Rel("device").URL())
	d.Set("externalips", resource.Rel("externalips").URL())
	d.Set("externalnetworks", resource.Rel("externalnetworks").URL())
	d.Set("network", vdcNetwork(resource))
	d.Set("privatenetworks", resource.Rel("privatenetworks").Href)
	d.Set("topurchase", resource.Rel("topurchase").Href)
	d.Set("purchased", resource.Rel("purchased").Href)
	d.Set("templates", resource.Rel("templates").Href)
	d.Set("tiers", resource.Rel("tiers").Href)
	purchaseIPs(resource, d.set("publicips"))
	// Default private network
	vdc := resource.(*abiquo.VirtualDatacenter)
	d.Set("net_name", vdc.Network.Name)
	d.Set("net_mask", vdc.Network.Mask)
	d.Set("net_gateway", vdc.Network.Gateway)
	d.Set("net_dns1", vdc.Network.DNS1)
	d.Set("net_dns2", vdc.Network.DNS2)
	d.Set("net_suffix", vdc.Network.Suffix)
	return
}

func vdcUpdate(d *resourceData, resource core.Resource) (err error) {
	if d.HasChange("publicips") {
		if err = purchaseIPs(resource, d.set("publicips")); err == nil {
			err = releaseIPs(resource, d.set("publicips"))
		}
	}
	return
}

func vdcRead(d *resourceData, resource core.Resource) (err error) {
	virtualdatacenter := resource.(*abiquo.VirtualDatacenter)
	// publicips
	publicips := schema.NewSet(schema.HashString, nil)
	virtualdatacenter.Rel("purchased").Collection(nil).List().Map(func(resource core.Resource) {
		publicips.Add(resource.Link().Title)
	})

	d.Set("name", virtualdatacenter.Name)
	d.Set("publicips", publicips)
	d.Set("cpuhard", virtualdatacenter.CPUHard)
	d.Set("cpusoft", virtualdatacenter.CPUSoft)
	d.Set("diskhard", virtualdatacenter.DiskHard)
	d.Set("disksoft", virtualdatacenter.DiskSoft)
	d.Set("publichard", virtualdatacenter.PublicHard)
	d.Set("publicsoft", virtualdatacenter.PublicSoft)
	d.Set("ramhard", virtualdatacenter.RAMHard)
	d.Set("ramsoft", virtualdatacenter.RAMSoft)
	d.Set("storagehard", virtualdatacenter.StorageHard)
	d.Set("storagesoft", virtualdatacenter.StorageSoft)
	d.Set("vlansoft", virtualdatacenter.VLANSoft)
	d.Set("vlanhard", virtualdatacenter.VLANHard)
	return
}

var virtualdatacenter = &description{
	media:    "virtualdatacenter",
	dto:      vdcNew,
	read:     vdcRead,
	update:   vdcUpdate,
	create:   vdcCreate,
	endpoint: endpointConst("cloud/virtualdatacenters"),
	Resource: &schema.Resource{Schema: vdcSchema},
}
