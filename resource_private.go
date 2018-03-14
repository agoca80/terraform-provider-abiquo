package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var privateSchema = map[string]*schema.Schema{
	"address": &schema.Schema{
		ForceNew: true,
		Required: true,
		Type:     schema.TypeString,
	},
	"mask": &schema.Schema{
		ForceNew: true,
		Required: true,
		Type:     schema.TypeInt,
	},
	"name": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"gateway": &schema.Schema{
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateIP,
	},
	"dns1": &schema.Schema{
		Optional:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateIP,
	},
	"dns2": &schema.Schema{
		Optional:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateIP,
	},
	"suffix": &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	},
	"virtualdatacenter": &schema.Schema{
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
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
