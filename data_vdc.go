package main

import (
	"errors"
	"fmt"
	"net/url"
	"path"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var vdcDataSource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name": Required().String(),
	},
	Read: dataVDCRead,
}

func dataVDCRead(d *schema.ResourceData, meta interface{}) (err error) {
	enterprise := meta.(*provider).Enterprise()
	if enterprise == nil {
		return errors.New("The user enterprise was not found")
	}

	id := path.Base(enterprise.URL())
	query := url.Values{"enterprise": {id}}
	finder := func(r core.Resource) bool {
		return r.(*abiquo.VirtualDatacenter).Name == d.Get("name").(string)
	}
	if vdc := enterprise.VirtualDatacenters(query).Find(finder); vdc != nil {
		d.SetId(vdc.URL())
		return
	}

	return fmt.Errorf("vdc %q was not found", d.Get("name"))
}
