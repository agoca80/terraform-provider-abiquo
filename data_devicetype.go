package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"

	"github.com/hashicorp/terraform/helper/schema"
)

var deviceTypeDataSchema = map[string]*schema.Schema{
	"name": Required().String(),
}

var deviceTypeDataSource = &schema.Resource{
	Schema: deviceTypeDataSchema,
	Read:   deviceTypeDataRead,
}

func deviceTypeDataRead(d *schema.ResourceData, meta interface{}) (err error) {
	name := d.Get("name").(string)
	finder := func(r core.Resource) bool {
		return r.(*abiquo.DeviceType).Name == name
	}
	resource := abiquo.DeviceTypes(nil).Find(finder)
	if resource == nil {
		return fmt.Errorf("device type %q does not exist", name)
	}
	d.SetId(resource.URL())
	return
}
