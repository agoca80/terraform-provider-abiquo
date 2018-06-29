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
	"virtualdatacenter": attribute(required, vdc),
}

func vappDataRead(d *schema.ResourceData, meta interface{}) (err error) {
	href := d.Get("virtualdatacenter").(string)
	vdc := core.NewLinker(href, "virtualdatacenter").Walk()
	if vdc == nil {
		return fmt.Errorf("virtualdatacenter %q not found", href)
	}

	name := d.Get("name").(string)
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
