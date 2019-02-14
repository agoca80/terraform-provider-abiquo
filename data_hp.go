package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var hpDataSchema = map[string]*schema.Schema{
	"name":     attribute(required, text),
	"location": attribute(required, link("location")),
}

func hpFind(d *resourceData) (err error) {
	name := d.string("name")
	href := d.string("location") + "/hardwareprofiles"
	hardwareprofiles := linkType(href, "hardwareprofiles").Collection(nil)
	hardwareprofile := hardwareprofiles.Find(func(r core.Resource) bool {
		return r.(*abiquo.HardwareProfile).Name == name
	})
	if hardwareprofile == nil {
		return fmt.Errorf("hwprofile %q does not exist in %q", name, href)
	}

	d.SetId(hardwareprofile.URL())
	return
}
