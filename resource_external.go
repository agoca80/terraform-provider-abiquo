package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var externalSchema = map[string]*schema.Schema{
	"datacenter":         endpoint("datacenter"),
	"address":            attribute(required, forceNew, ip),
	"tag":                attribute(required, forceNew, positive),
	"mask":               attribute(required, forceNew, positive),
	"name":               attribute(required, text),
	"gateway":            attribute(required, ip),
	"dns1":               attribute(optional, ip),
	"dns2":               attribute(optional, ip),
	"suffix":             attribute(optional, text),
	"networkservicetype": attribute(required, forceNew, href),
	"enterprise":         attribute(required, forceNew, link("enterprise")),
}

func externalNew(d *resourceData) core.Resource {
	network := networkNew(d)
	network.Type = "EXTERNAL"
	network.Tag = d.integer("tag")
	network.DTO = core.NewDTO(
		d.link("enterprise"),
		d.link("networkservicetype"),
	)
	return network
}

func externalRead(d *resourceData, resource core.Resource) (e error) {
	network := resource.(*abiquo.Network)
	networkRead(d, network)
	d.Set("enterprise", network.Rel("enterprise").URL())
	d.Set("nst", network.Rel("networkservicetype").URL())
	// d.Set("datacenter", network.Rel("datacenter").URL())
	return
}

var external = &description{
	name:     "external",
	Resource: &schema.Resource{Schema: externalSchema},
	dto:      externalNew,
	endpoint: endpointPath("datacenter", "/network"),
	media:    "vlan",
	read:     externalRead,
}
