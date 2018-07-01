package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var publicSchema = map[string]*schema.Schema{
	"address": attribute(required, ip, forceNew),
	"tag":     attribute(required, natural, forceNew),
	"mask":    attribute(required, natural, forceNew),
	"name":    attribute(required, text),
	"gateway": attribute(required, ip),
	"dns1":    attribute(optional, ip),
	"dns2":    attribute(optional, ip),
	"suffix":  attribute(optional, text),
	// Links
	"networkservicetype": attribute(required, href, forceNew),
	"datacenter":         attribute(required, datacenter, forceNew),
}

func publicNew(d *resourceData) core.Resource {
	public := networkNew(d)
	public.Type = "EXTERNAL"
	public.Tag = d.integer("tag")
	public.DTO = core.NewDTO(
		d.link("networkservicetype"),
	)
	return public
}

func publicEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("datacenter")+"/network", "vlan")
}

// PENDING: Public networks are not supossed to change, but ...
func publicRead(d *resourceData, resource core.Resource) (e error) {
	network := resource.(*abiquo.Network)
	networkRead(d, network)
	d.Set("networkservicetype", network.Rel("networkservicetype").URL())
	return
}
