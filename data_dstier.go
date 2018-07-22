package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"

	"github.com/hashicorp/terraform/helper/schema"
)

var dstierDataSchema = map[string]*schema.Schema{
	"name":       attribute(required, text),
	"datacenter": attribute(required, link("datacenter")),
}

func dstierFind(d *resourceData) (err error) {
	href := d.string("datacenter")
	datacenter := core.NewLinker(href, "datacenter").Walk()
	if datacenter == nil {
		return fmt.Errorf("datacenter not found: %q", href)
	}

	name := d.string("name")
	dstiers := datacenter.Rel("datastoretiers").Collection(nil)
	dstier := dstiers.Find(func(r core.Resource) bool {
		return r.(*abiquo.DatastoreTier).Name == name
	})
	if dstier == nil {
		return fmt.Errorf("datastore tier not found: %q", name)
	}

	d.SetId(dstier.URL())
	return
}
