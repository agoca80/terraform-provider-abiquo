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

func dstierDataRead(d *schema.ResourceData, meta interface{}) (err error) {
	href := d.Get("datacenter").(string)
	datacenter := core.NewLinker(href, "datacenter").Walk()
	if datacenter == nil {
		return fmt.Errorf("datacenter not found: %q", href)
	}

	name := d.Get("name").(string)
	dstier := datacenter.Rel("datastoretiers").Collection(nil).Find(func(r core.Resource) bool {
		return r.(*abiquo.DatastoreTier).Name == name
	})
	if dstier == nil {
		return fmt.Errorf("datastore tier not found: %q", name)
	}

	d.SetId(dstier.URL())
	return
}
