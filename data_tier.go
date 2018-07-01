package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"

	"github.com/hashicorp/terraform/helper/schema"
)

var tierDataSchema = map[string]*schema.Schema{
	"name":       attribute(required, text),
	"datacenter": attribute(required, link("datacenter")),
}

func tierDataRead(d *schema.ResourceData, meta interface{}) (err error) {
	href := d.Get("datacenter").(string)
	datacenter := core.NewLinker(href, "datacenter").Walk()
	if datacenter == nil {
		return fmt.Errorf("datacenter not found: %q", href)
	}

	name := d.Get("name").(string)
	tier := datacenter.Rel("tiers").Collection(nil).Find(func(r core.Resource) bool {
		return r.(*abiquo.Tier).Name == name
	})
	if tier == nil {
		return fmt.Errorf("datastore tier not found: %q", name)
	}

	d.SetId(tier.URL())
	return
}
