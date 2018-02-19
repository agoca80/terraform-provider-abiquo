package main

import (
	"github.com/abiquo/opal/abiquo"
	"github.com/abiquo/opal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var rackSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"datacenter": &schema.Schema{
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
}

var rackResource = &schema.Resource{
	Schema: rackSchema,
	Delete: resourceDelete,
	Exists: resourceExists("rack"),
	Create: resourceCreate(rackNew, nil, rackRead, rackEndpoint),
	Update: resourceUpdate(rackNew, "rack"),
	Read:   resourceRead(rackNew, rackRead, "rack"),
}

func rackNew(d *resourceData) core.Resource {
	return &abiquo.Rack{
		Name: d.string("name"),
	}
}

func rackEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("datacenter")+"/racks", "rack")
}

func rackRead(d *resourceData, resource core.Resource) (err error) {
	rack := resource.(*abiquo.Rack)
	d.Set("name", rack.Name)
	return
}
