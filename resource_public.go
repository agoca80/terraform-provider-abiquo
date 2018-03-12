package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var publicResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"address":            Required().Renew().String(),
		"tag":                Required().Renew().Number(),
		"mask":               Required().Renew().Number(),
		"name":               Required().String(),
		"gateway":            Required().IP(),
		"dns1":               Optional().IP(),
		"dns2":               Optional().IP(),
		"suffix":             Optional().String(),
		"networkservicetype": Required().Renew().Link(),
		"datacenter":         Required().Renew().Link(),
	},
	Delete: resourceDelete,
	Update: resourceUpdate(publicNew, nil, "vlan"),
	Create: resourceCreate(publicNew, nil, publicRead, publicEndpoint),
	Read:   resourceRead(publicNew, publicRead, "vlan"),
}

func publicNew(d *resourceData) core.Resource {
	public := networkNew(d)
	public.TypeNet = "EXTERNAL"
	public.Tag = d.int("tag")
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
