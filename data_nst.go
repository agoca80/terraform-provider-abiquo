package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"

	"github.com/hashicorp/terraform/helper/schema"
)

var nstDataSchema = map[string]*schema.Schema{
	"name":       attribute(required, text),
	"datacenter": attribute(required, link("datacenter")),
}

func nstDataRead(d *schema.ResourceData, meta interface{}) (err error) {
	href := d.Get("datacenter").(string)
	endpoint := core.NewLinker(href, "datacenter")
	resource := endpoint.Walk()
	if resource == nil {
		return fmt.Errorf("datacenter not found: %q", href)
	}

	name := d.Get("name").(string)
	datacenter := resource.(*abiquo.Datacenter)
	resource = datacenter.Rel("networkservicetypes").Collection(nil).Find(func(r core.Resource) bool {
		return r.(*abiquo.NetworkServiceType).Name == name
	})
	if resource == nil {
		return fmt.Errorf("network service type not found: %q", name)
	}

	d.SetId(resource.URL())
	return
}
