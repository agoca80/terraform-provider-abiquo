package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var externalSchema = map[string]*schema.Schema{
	"address": &schema.Schema{
		ForceNew: true,
		Required: true,
		Type:     schema.TypeString,
	},
	"tag": &schema.Schema{
		ForceNew: true,
		Required: true,
		Type:     schema.TypeInt,
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
	"enterprise": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
	"networkservicetype": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
	"datacenter": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
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
