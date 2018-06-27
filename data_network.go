package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"

	"github.com/hashicorp/terraform/helper/schema"
)

var networkDataSchema = map[string]*schema.Schema{
	"ips": &schema.Schema{
		Computed: true,
		Type:     schema.TypeString,
	},
	"name": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"location": &schema.Schema{
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
}

func networkDataRead(d *schema.ResourceData, meta interface{}) (err error) {
	name := d.Get("name").(string)
	href := d.Get("location").(string)
	networks := core.NewLinkType(href, "vlans").Collection(nil)
	network := networks.Find(func(r core.Resource) bool {
		return r.(*abiquo.Network).Name == name
	})
	if network == nil {
		return fmt.Errorf("network %q not found", name)
	}

	d.SetId(network.URL())
	d.Set("ips", network.Rel("ips").Href)
	return
}
