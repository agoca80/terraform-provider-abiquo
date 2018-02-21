package main

import (
	"fmt"

	"github.com/abiquo/opal/abiquo"
	"github.com/abiquo/opal/core"

	"github.com/hashicorp/terraform/helper/schema"
)

var dstierDataSchema = map[string]*schema.Schema{
	"name":       Required().String(),
	"datacenter": Required().ValidateURL(),
}

var dstierDataSource = &schema.Resource{
	Schema: dstierDataSchema,
	Read:   dstierDataRead,
}

func dstierDataRead(d *schema.ResourceData, meta interface{}) (err error) {
	href := d.Get("datacenter").(string)
	resource := core.NewLinker(href, "datacenter").Walk()
	if resource == nil {
		return fmt.Errorf("datacenter not found: %q", href)
	}

	name := d.Get("name").(string)
	datacenter := resource.(*abiquo.Datacenter)
	resource = datacenter.Rel("datastoretiers").Collection(nil).Find(func(r core.Resource) bool {
		return r.(*abiquo.DatastoreTier).Name == name
	})
	if resource == nil {
		return fmt.Errorf("datastore tier not found: %q", name)
	}

	d.SetId(resource.URL())
	return
}
