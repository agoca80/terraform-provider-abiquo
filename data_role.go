package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var roleDataSchema = map[string]*schema.Schema{
	"name": attribute(required, text),
}

func roleFind(d *resourceData) (err error) {
	name := d.string("name")
	resource := abiquo.Roles(nil).Find(func(r core.Resource) bool {
		return r.(*abiquo.Role).Name == name
	})
	if resource == nil {
		return fmt.Errorf("role %q was not found", d.Get("name"))
	}
	d.SetId(resource.URL())
	return
}
