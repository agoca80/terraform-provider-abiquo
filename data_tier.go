package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"

	"github.com/hashicorp/terraform/helper/schema"
)

var tierDataSchema = map[string]*schema.Schema{
	"name":     attribute(required, text),
	"location": attribute(required, href),
}

func tierFind(d *resourceData) (err error) {
	name := d.string("name")
	href := d.string("location")
	endpoint := linkType(href, "tiers")
	tier := endpoint.Collection(nil).Find(func(r core.Resource) bool {
		return r.(*abiquo.Tier).Name == name
	})
	if tier == nil {
		return fmt.Errorf("tier not found: %q", name)
	}

	d.SetId(tier.URL())
	return
}
