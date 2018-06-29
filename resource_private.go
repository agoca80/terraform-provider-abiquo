package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var privateSchema = map[string]*schema.Schema{
	"address": attribute(required, ip, forceNew),
	"mask":    attribute(required, natural, forceNew),
	"name":    attribute(required, text),
	"gateway": attribute(required, ip),
	"dns1":    attribute(optional, ip),
	"dns2":    attribute(optional, ip),
	"suffix":  attribute(optional, text),
	// Links
	"virtualdatacenter": attribute(required, vdc, forceNew),
}

func privateEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("virtualdatacenter")+"/privatenetworks", "vlan")
}

func privateNew(d *resourceData) core.Resource {
	private := networkNew(d)
	private.Type = "INTERNAL"
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
