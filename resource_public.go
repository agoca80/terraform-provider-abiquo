package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var publicSchema = map[string]*schema.Schema{
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

var publicResource = &schema.Resource{
	Schema: publicSchema,
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
