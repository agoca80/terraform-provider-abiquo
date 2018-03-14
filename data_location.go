package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var locationDataSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
}

var locationDataSource = &schema.Resource{
	Schema: locationDataSchema,
	Read:   locationRead,
}

func locationRead(d *schema.ResourceData, meta interface{}) (err error) {
	finder := func(r core.Resource) bool {
		return r.(*abiquo.Datacenter).Name == d.Get("name").(string)
	}
	location := abiquo.PrivateLocations(nil).Find(finder)
	if location == nil {
		return fmt.Errorf("Location %q does not exist", d.Get("name"))
	}

	d.SetId(location.Rel("location").Href)
	return
}
