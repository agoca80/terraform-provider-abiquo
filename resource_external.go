package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var externalResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"address":            Required().Renew().String(),
		"tag":                Required().Renew().Number(),
		"mask":               Required().Renew().Number(),
		"name":               Required().String(),
		"gateway":            Required().IP(),
		"dns1":               Optional().IP(),
		"dns2":               Optional().IP(),
		"suffix":             Optional().String(),
		"enterprise":         Required().Renew().Link(),
		"networkservicetype": Required().Renew().Link(),
		"datacenter":         Required().Renew().Link(),
	},
	Delete: resourceDelete,
	Exists: resourceExists("vlan"),
	Update: resourceUpdate(externalNew, nil, "vlan"),
	Create: resourceCreate(externalNew, nil, externalRead, externalEndpoint),
	Read:   resourceRead(externalNew, externalRead, "vlan"),
}

func externalNew(d *resourceData) core.Resource {
	network := networkNew(d)
	network.TypeNet = "EXTERNAL"
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
