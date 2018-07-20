package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var dstierSchema = map[string]*schema.Schema{
	"datacenter":  endpoint("datacenter"),
	"description": attribute(required, text),
	"enabled":     attribute(required, boolean),
	"name":        attribute(required, text),
	"policy":      attribute(required, label([]string{"PERFORMANCE", "PROGRESSIVE"})),
}

func dstierDTO(d *resourceData) core.Resource {
	return &abiquo.DatastoreTier{
		Description: d.string("description"),
		Enabled:     d.boolean("enabled"),
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

var resourceDstier = &schema.Resource{
	Schema: dstierSchema,
	Delete: resourceDelete,
	Exists: resourceExists("datastoretier"),
	Create: resourceCreate(dstierDTO, nil, dstierRead, dstierEndpoint),
	Update: resourceUpdate(dstierDTO, nil, "datastoretier"),
	Read:   resourceRead(dstierDTO, dstierRead, "datastoretier"),
}
