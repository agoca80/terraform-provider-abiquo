package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var templateDataSchema = map[string]*schema.Schema{
	"templates": attribute(required, link("templates")),
	"name":      attribute(required, text),
}

func templateFind(d *resourceData) (err error) {
	name := d.string("name")
	templates := d.string("templates")
	endpoint := linkType(templates, "virtualmachinetemplates")
	template := endpoint.Collection(nil).Find(func(r core.Resource) bool {
		t := r.(*abiquo.VirtualMachineTemplate)
		return t.Name == name && t.State != "UNAVAILABLE"
	})
	if template == nil {
		return fmt.Errorf("template %q was not found", d.Get("name"))
	}
	d.SetId(template.URL())
	return
}
