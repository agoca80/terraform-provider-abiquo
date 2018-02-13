package main

import (
	"fmt"
	"sync"

	"github.com/abiquo/opal/abiquo"
	"github.com/hashicorp/terraform/helper/schema"
)

var datacenters = struct {
	sync.Once
	datacenter map[string]*abiquo.Datacenter
}{}

var datacenterDataSource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name": Required().String(),
	},
	Read: datacenterRead,
}

func datacenterRead(d *schema.ResourceData, meta interface{}) (err error) {
	datacenter := datacenterGet(d.Get("name").(string))
	if datacenter == nil {
		return fmt.Errorf("Datacenter %v does not exist", d.Get("name"))
	}
	d.SetId(datacenter.URL())
	return
}

func datacenterGet(name string) *abiquo.Datacenter {
	datacenters.Do(func() {
		datacenters.datacenter = make(map[string]*abiquo.Datacenter)
		for _, d := range abiquo.Datacenters(nil).List() {
			datacenter := d.(*abiquo.Datacenter)
			datacenters.datacenter[datacenter.Name] = datacenter
		}
	})
	return datacenters.datacenter[name]
}

func datacenterID(name interface{}) int {
	return datacenterGet(name.(string)).ID
}

func datacenterName(id int) string {
	for name, dc := range datacenters.datacenter {
		if id == dc.ID {
			return name
		}
	}
	return ""
}
