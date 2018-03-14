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

var vdcDataSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
}

var vdcDataSource = &schema.Resource{
	Schema: vdcDataSchema,
	Read:   dataVDCRead,
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
	vdcs := enterprise.Rel("cloud/virtualdatacenters").Collection(query)
	if vdc := vdcs.Find(finder); vdc != nil {
		d.SetId(vdc.URL())
		return
	}

	return fmt.Errorf("vdc %q was not found", d.Get("name"))
}
