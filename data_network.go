package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"

	"github.com/hashicorp/terraform/helper/schema"
)

var networkDataSchema = map[string]*schema.Schema{
	"ips":      attribute(computed, text),
	"name":     attribute(required, text),
	"location": attribute(required, link("location")),
}

func networkFind(d *resourceData) (err error) {
	name := d.string("name")
	href := d.string("location")
	networks := core.NewLinkType(href, "vlans").Collection(nil)
	network := networks.Find(func(r core.Resource) bool {
		return r.(*abiquo.Network).Name == name
	})
	if network == nil {
		return fmt.Errorf("network %q not found", name)
	}

	d.SetId(network.URL())
	d.Set("ips", network.Rel("ips").Href)
	return
}
