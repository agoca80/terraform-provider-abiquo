package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var roleDataSchema = map[string]*schema.Schema{
	"name": Required().String(),
}

var roleDataSource = &schema.Resource{
	Schema: roleDataSchema,
	Read:   roleDataRead,
}

func roleDataRead(d *schema.ResourceData, meta interface{}) (err error) {
	name := d.Get("name").(string)
	finder := func(r core.Resource) bool {
		return r.(*abiquo.Role).Name == name
	}
	resource := abiquo.Roles(nil).Find(finder)
	if resource == nil {
		return fmt.Errorf("role %q was not found", d.Get("name"))
	}
	d.SetId(resource.URL())
	return
}
