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

func nstFind(d *resourceData) (err error) {
	href := d.string("datacenter")
	endpoint := core.NewLinker(href, "datacenter")
	resource, err := endpoint.Walk()
	if err != nil {
		return
	}

	name := d.string("name")
	datacenter := resource.(*abiquo.Datacenter)
	nsts := datacenter.Rel("networkservicetypes").Collection(nil)
	nst := nsts.Find(func(r core.Resource) bool {
		return r.(*abiquo.NetworkServiceType).Name == name
	})
	if nst == nil {
		return fmt.Errorf("network service type not found: %q", name)
	}

	d.SetId(nst.URL())
	return
}
