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

func privateDTO(d *resourceData) core.Resource {
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

var private = &description{
	name:     "private",
	dto:      privateDTO,
	endpoint: endpointPath("virtualdatacenter", "/privatenetworks"),
	media:    "vlan",
	read:     privateRead,
	Resource: &schema.Resource{Schema: privateSchema},
}
