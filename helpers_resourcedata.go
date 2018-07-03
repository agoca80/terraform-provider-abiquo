package main

import (
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

type resourceData struct {
	*core.Link
	*schema.ResourceData
}

func newResourceData(d *schema.ResourceData, media string) *resourceData {
	return &resourceData{
		Link:         core.NewLinkType(d.Id(), media),
		ResourceData: d,
	}
}

func (d *resourceData) slice(name string) (slice []interface{}) {
	if i, ok := d.GetOk(name); ok {
		slice = i.([]interface{})
	}
	return
}

func (d *resourceData) dict(name string) (m map[string]interface{}) {
	if i, ok := d.GetOk(name); ok {
		m = i.(map[string]interface{})
	}
	return
}

func (d *resourceData) set(name string) (s *schema.Set) {
	if i, ok := d.GetOk(name); ok {
		s = i.(*schema.Set)
	}
	return
}

func (d *resourceData) link(name string) (link *core.Link) {
	if _, ok := d.GetOk(name); ok {
		link = core.NewLinkType(d.string(name), name).SetRel(name)
	}
	return
}

func (d *resourceData) string(name string) string {
	return d.ResourceData.Get(name).(string)
}

func (d *resourceData) integer(name string) (val int) {
	if i, ok := d.GetOk(name); ok {
		val = i.(int)
	}
	return
}

func (d *resourceData) boolean(name string) (val bool) {
	if i, ok := d.GetOk(name); ok {
		val = i.(bool)
	}
	return
}
