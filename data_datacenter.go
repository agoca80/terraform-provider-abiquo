package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var datacenterDataSchema = map[string]*schema.Schema{
	"name":    attribute(required, text),
	"network": attribute(computed, text),
	"tiers":   attribute(computed, text),
}

func datacenterDataRead(d *schema.ResourceData, meta interface{}) (err error) {
	name := d.Get("name").(string)
	datacenter := abiquo.Datacenters(nil).Find(func(r core.Resource) bool {
		return r.(*abiquo.Datacenter).Name == name
	})
	if datacenter == nil {
		return fmt.Errorf("Datacenter %v does not exist", d.Get("name"))
	}
	d.SetId(datacenter.URL())
	d.Set("network", datacenter.Rel("network").Href)
	d.Set("tiers", datacenter.Rel("tiers").Href)
	return
}
