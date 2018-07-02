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
	"name":    attribute(required, text),
	"tiers":   attribute(computed, text),
	"network": attribute(computed, text),
}

func vdcNetwork(r core.Resource) *core.Link {
	vdc := r.(*abiquo.VirtualDatacenter)
	for _, l := range vdc.Links {
		if l.Title == "default_private_network" {
			return l
		}
	}
	return nil
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
	vdc := enterprise.Rel("cloud/virtualdatacenters").Collection(query).Find(finder)
	if vdc == nil {
		return fmt.Errorf("vdc %q was not found", d.Get("name"))
	}

	d.SetId(vdc.URL())
	d.Set("tiers", vdc.Rel("tiers").Href)
	d.Set("network", vdcNetwork(vdc).Href)
	return
}
