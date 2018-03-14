package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var dstierSchema = map[string]*schema.Schema{
	"datacenter": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
	"description": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"enabled": &schema.Schema{
		Required: true,
		Type:     schema.TypeBool,
	},
	"name": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"policy": &schema.Schema{
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateString([]string{"PERFORMANCE", "PROGRESSIVE"}),
	},
}

func dstierDTO(d *resourceData) core.Resource {
	return &abiquo.DatastoreTier{
		Description: d.string("description"),
		Enabled:     d.bool("enabled"),
		Name:        d.string("name"),
		Policy:      d.string("policy"),
	}
}

func dstierEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("datacenter")+"/datastoretiers", "datastoretier")
}

func dstierRead(d *resourceData, resource core.Resource) (err error) {
	dstier := resource.(*abiquo.DatastoreTier)
	d.Set("description", dstier.Description)
	d.Set("enabled", dstier.Enabled)
	d.Set("name", dstier.Name)
	d.Set("policy", dstier.Policy)
	d.Set("datacenter", dstier.Rel("datacenter").URL())
	return
}
