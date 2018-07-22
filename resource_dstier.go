package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var datastoretierSchema = map[string]*schema.Schema{
	"datacenter":  endpoint("datacenter"),
	"description": attribute(required, text),
	"enabled":     attribute(required, boolean),
	"name":        attribute(required, text),
	"policy":      attribute(required, label([]string{"PERFORMANCE", "PROGRESSIVE"})),
}

func datastoretierDTO(d *resourceData) core.Resource {
	return &abiquo.DatastoreTier{
		Description: d.string("description"),
		Enabled:     d.boolean("enabled"),
		Name:        d.string("name"),
		Policy:      d.string("policy"),
	}
}

func datastoretierRead(d *resourceData, resource core.Resource) (err error) {
	datastoretier := resource.(*abiquo.DatastoreTier)
	d.Set("description", datastoretier.Description)
	d.Set("enabled", datastoretier.Enabled)
	d.Set("name", datastoretier.Name)
	d.Set("policy", datastoretier.Policy)
	d.Set("datacenter", datastoretier.Rel("datacenter").URL())
	return
}

var datastoretier = &description{
	media:    "datastoretier",
	dto:      datastoretierDTO,
	read:     datastoretierRead,
	endpoint: endpointPath("datacenter", "/datastoretiers"),
	Resource: &schema.Resource{Schema: datastoretierSchema},
}
