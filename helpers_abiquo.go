package main

import (
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/hilmapstructure"
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
