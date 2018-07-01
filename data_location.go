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
	datacenters := abiquo.PrivateLocations(nil)
	if location = datacenters.Find(func(r core.Resource) bool {
		return r.(*abiquo.Datacenter).Name == name
	}); location != nil {
		return
	}

	regions := abiquo.PublicLocations(nil)
	location = regions.Find(func(r core.Resource) bool {
		return r.(*abiquo.Location).Name == name
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
