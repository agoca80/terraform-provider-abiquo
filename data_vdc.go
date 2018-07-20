package main

import (
	"fmt"
	"net/url"
	"path"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var vdcDataSchema = map[string]*schema.Schema{
	"name":      attribute(required, text),
	"tiers":     attribute(computed, text),
	"location":  attribute(computed, text),
	"network":   attribute(computed, text),
	"templates": attribute(computed, text),
}

func vdcNetwork(r core.Resource) string {
	vdc := r.(*abiquo.VirtualDatacenter)
	network := vdc.Links.Find(func(l *core.Link) bool {
		return l.Title == "default_private_network"
	})
	return network.URL()
}

func dataVDCRead(d *schema.ResourceData, meta interface{}) (err error) {
	enterprise := meta.(*provider).Enterprise()
	id := path.Base(enterprise.URL())
	query := url.Values{"enterprise": {id}}
	vdcs := enterprise.Rel("cloud/virtualdatacenters").Collection(query)
	vdc := vdcs.Find(func(r core.Resource) bool {
		return r.(*abiquo.VirtualDatacenter).Name == d.Get("name").(string)
	})
	if vdc == nil {
		return fmt.Errorf("vdc %q was not found", d.Get("name"))
	}

	d.SetId(vdc.URL())
	d.Set("tiers", vdc.Rel("tiers").Href)
	d.Set("network", vdcNetwork(vdc))
	d.Set("location", vdc.Rel("location").Href)
	d.Set("templates", vdc.Rel("templates").Href)
	return
}

var dataVdc = &schema.Resource{
	Schema: vdcDataSchema,
	Read:   dataVDCRead,
}
