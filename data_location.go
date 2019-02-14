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

func lFind(name string) (location core.Resource) {
	if location = abiquo.PublicLocations().Collection(nil).Find(func(r core.Resource) bool {
		return r.(*abiquo.Location).Name == name
	}); location != nil {
		return
	}

	location = abiquo.PrivateLocations().Collection(nil).Find(func(r core.Resource) bool {
		return r.(*abiquo.Datacenter).Name == name
	})

	return
}

func locationFind(d *resourceData) (err error) {
	if location := lFind(d.string("name")); location != nil {
		d.SetId(location.Rel("location").Href)
		return
	}
	return fmt.Errorf("Location %q does not exist", d.Get("name"))
}
