package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var externalSchema = map[string]*schema.Schema{
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
	"enterprise":         attribute(required, enterprise, forceNew),
}

func externalNew(d *resourceData) core.Resource {
	network := networkNew(d)
	network.Type = "EXTERNAL"
	network.Tag = d.int("tag")
	network.DTO = core.NewDTO(
		d.link("enterprise"),
		d.link("networkservicetype"),
	)
	return network
}

func externalEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("datacenter")+"/network", "vlan")
}

func externalRead(d *resourceData, resource core.Resource) (e error) {
	network := resource.(*abiquo.Network)
	networkRead(d, network)
	d.Set("enterprise", network.Rel("enterprise").URL())
	d.Set("nst", network.Rel("networkservicetype").URL())
	// d.Set("datacenter", network.Rel("datacenter").URL())
	return
}
