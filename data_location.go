package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var locationDataSchema = map[string]*schema.Schema{
	"name": attribute(required, text),
}

func locationFind(name string) (location core.Resource) {
	if location = abiquo.PublicLocations(nil).Find(func(r core.Resource) bool {
		return r.(*abiquo.Location).Name == name
	}); location != nil {
		return
	}

	location = abiquo.PrivateLocations(nil).Find(func(r core.Resource) bool {
		return r.(*abiquo.Datacenter).Name == name
	})

	return
}

func locationRead(d *schema.ResourceData, meta interface{}) (err error) {
	if location := locationFind(d.Get("name").(string)); location != nil {
		d.SetId(location.Rel("location").Href)
		return
	}
	return fmt.Errorf("Location %q does not exist", d.Get("name"))
}

var dataLocation = &schema.Resource{
	Schema: locationDataSchema,
	Read:   locationRead,
}
