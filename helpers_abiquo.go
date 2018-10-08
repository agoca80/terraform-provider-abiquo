package main

import (
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/hilmapstructure"
	"github.com/hashicorp/terraform/helper/schema"
)

func mapDecoder(m interface{}, i interface{}) interface{} {
	if err := hilmapstructure.WeakDecode(m.(map[string]interface{}), i); err != nil {
		panic("mapDecoder: error decoding map")
	}
	return i
}

func mapHrefs(links core.Links) (hrefs []interface{}) {
	for _, l := range links {
		hrefs = append(hrefs, l.Href)
	}
	return
}

type method func(*resourceData) error

func data(s map[string]*schema.Schema, find method) *schema.Resource {
	return &schema.Resource{
		Schema: s,
		Read: func(rd *schema.ResourceData, _ interface{}) error {
			return find(newData(rd))
		},
	}
}
