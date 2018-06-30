package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"

	"github.com/hashicorp/terraform/helper/schema"
)

var deviceTypeDataSchema = map[string]*schema.Schema{
	"name": attribute(required, text),
}

func deviceTypeDataRead(d *schema.ResourceData, meta interface{}) (err error) {
	name := d.Get("name").(string)
	resource := abiquo.DeviceTypes(nil).Find(func(r core.Resource) bool {
		return r.(*abiquo.DeviceType).Name == name
	})
	if resource == nil {
		return fmt.Errorf("device type %q does not exist", name)
	}
	d.SetId(resource.URL())
	return
}
