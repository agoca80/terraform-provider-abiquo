package main

import (
	"github.com/abiquo/opal/abiquo"
	"github.com/abiquo/opal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var privateResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"address":           Required().Renew().String(),
		"mask":              Required().Renew().Number(),
		"name":              Required().String(),
		"gateway":           Required().IP(),
		"dns1":              Optional().IP(),
		"dns2":              Optional().IP(),
		"suffix":            Optional().String(),
		"ips":               optional(fieldMap(fieldIP())),
		"virtualdatacenter": Required().Link(),
	},
	Delete: resourceDelete,
	Update: resourceUpdate(privateNew, "vlan"),
	Create: resourceCreate(privateNew, privateEndpoint),
	Read:   resourceRead(privateNew, privateRead, "vlan"),
}

func privateEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("virtualdatacenter")+"/privatenetworks", "vlan")
}

func privateNew(d *resourceData) core.Resource {
	private := networkNew(d)
	private.TypeNet = "INTERNAL"
	private.DTO = core.NewDTO(
		d.link("virtualdatacenter"),
	)
	return private
}

func privateRead(d *resourceData, resource core.Resource) (e error) {
	network := resource.(*abiquo.Network)
	networkRead(d, network)
	d.Set("virtualdatacenter", network.Rel("virtualdatacenter").URL())
	return
}
