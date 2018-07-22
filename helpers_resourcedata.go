package main

import (
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

type resourceData struct {
	*core.Link
	*schema.ResourceData
}

func newData(rd *schema.ResourceData) (d *resourceData) {
	return &resourceData{
		ResourceData: rd,
	}
}

func newDataType(rd *schema.ResourceData, media string) (d *resourceData) {
	d = newData(rd)
	d.Link = core.NewLinkType(rd.Id(), media)
	return
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

func (d *resourceData) SetOk(name string, value interface{}) {
	if _, ok := d.GetOk(name); ok {
		d.Set(name, value)
	}
}

func (d *resourceData) link(name string) (link *core.Link) {
	if _, ok := d.GetOk(name); ok {
		link = core.NewLinkType(d.string(name), name).SetRel(name)
	}
	return
}

func (d *resourceData) string(name string) string {
	return d.Get(name).(string)
}

func (d *resourceData) integer(name string) (val int) {
	if i, ok := d.GetOk(name); ok {
		val = i.(int)
	}
	return
}

func (d *resourceData) float(name string) (val float64) {
	if i, ok := d.GetOk(name); ok {
		val = i.(float64)
	}
	return
}

func (d *resourceData) boolean(name string) (val bool) {
	if i, ok := d.GetOk(name); ok {
		val = i.(bool)
	}
	return
}
