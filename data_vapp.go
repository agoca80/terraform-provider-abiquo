package main

import (
	"fmt"
	"net/url"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"

	"github.com/hashicorp/terraform/helper/schema"
)

var vappDataSchema = map[string]*schema.Schema{
	"name":              attribute(required, text),
	"virtualdatacenter": attribute(required, link("virtualdatacenter")),
}

func vappFind(d *resourceData) (err error) {
	href := d.string("virtualdatacenter")
	vdc, err := linkType(href, "virtualdatacenter").Walk()
	if err != nil {
		return
	}

	name := d.string("name")
	query := url.Values{"has": {name}}
	vapps := vdc.Rel("virtualappliances").Collection(query)
	vapp := vapps.Find(func(r core.Resource) bool {
		return r.(*abiquo.VirtualAppliance).Name == name
	})
	if vapp == nil {
		return fmt.Errorf("virtual appliance %q not found", name)
	}

	d.SetId(vapp.URL())
	return
}
