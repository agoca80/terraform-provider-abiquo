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

func tierDataRead(d *schema.ResourceData, meta interface{}) (err error) {
	name := d.Get("name").(string)
	href := d.Get("location").(string)
	endpoint := core.NewLinkType(href, "tiers")
	tier := endpoint.Collection(nil).Find(func(r core.Resource) bool {
		return r.(*abiquo.Tier).Name == name
	})
	if tier == nil {
		return fmt.Errorf("tier not found: %q", name)
	}

	d.SetId(tier.URL())
	return
}

var dataTier = &schema.Resource{
	Schema: tierDataSchema,
	Read:   tierDataRead,
}
