package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var privateSchema = map[string]*schema.Schema{
	"virtualdatacenter": endpoint("virtualdatacenter"),
	"address":           attribute(required, forceNew, ip),
	"mask":              attribute(required, forceNew, natural),
	"name":              attribute(required, text),
	"gateway":           attribute(required, ip),
	"dns1":              attribute(optional, ip),
	"dns2":              attribute(optional, ip),
	"suffix":            attribute(optional, text),
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

var resourcePrivate = &schema.Resource{
	Schema: privateSchema,
	Delete: resourceDelete,
	Update: resourceUpdate(privateNew, nil, "vlan"),
	Create: resourceCreate(privateNew, nil, privateRead, privateEndpoint),
	Read:   resourceRead(privateNew, privateRead, "vlan"),
}
