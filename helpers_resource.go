package main

import (
	"fmt"

	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

type description struct {
	media    string
	name     string
	dto      func(*resourceData) core.Resource
	endpoint func(*resourceData) string
	create   func(*resourceData, core.Resource) error
	read     func(*resourceData, core.Resource) error
	update   func(*resourceData, core.Resource) error
	*schema.Resource
}

func (d *description) readFn() schema.ReadFunc {
	return func(rd *schema.ResourceData, _ interface{}) (err error) {
		data := newDataType(rd, d.media)
		resource := core.Factory(d.media)
		err = core.Read(data, resource)
		if err != nil {
			return fmt.Errorf("readFn: %v", err)
		}
		return d.read(data, resource)
	}
}

func (d *description) existsFn() schema.ExistsFunc {
	return func(rd *schema.ResourceData, _ interface{}) (ok bool, err error) {
		return newDataType(rd, d.media).Link.Exists()
	}
}

func (d *description) updateFn() schema.UpdateFunc {
	return func(rd *schema.ResourceData, m interface{}) (err error) {
		data := newDataType(rd, d.media)
		resource := d.dto(data)
		err = core.Update(data, resource)
		if err == nil && d.update != nil {
			err = d.update(data, resource)
		}
		return
	}
}

func (d *description) createFn() schema.CreateFunc {
	return func(rd *schema.ResourceData, _ interface{}) (err error) {
		data := newData(rd)
		resource := d.dto(data)
		if resource == nil {
			return fmt.Errorf("createFn: resource could not be created")
		}

		endpoint := core.NewLinker(d.endpoint(data), d.media)
		if err = core.Create(endpoint, resource); err != nil {
			return
		}
		data.SetId(resource.URL())
		if d.create != nil {
			err = d.create(data, resource)
		}

		if err == nil {
			err = d.read(data, resource)
		}

		return
	}
}

func resourceDelete(d *schema.ResourceData, m interface{}) (err error) {
	return core.Remove(newDataType(d, ""))
}

func endpointConst(href string) func(*resourceData) string {
	return func(data *resourceData) string {
		return href
	}
}

func endpointPath(base, path string) func(*resourceData) string {
	return func(data *resourceData) string {
		return data.string(base) + path
	}
}

func resourceDefinition(d *description) (r *schema.Resource) {
	r = d.Resource

	if r.Create == nil {
		r.Create = d.createFn()
	}

	if r.Exists == nil {
		r.Exists = d.existsFn()
	}

	if r.Read == nil {
		r.Read = d.readFn()
	}

	if r.Delete == nil {
		r.Delete = resourceDelete
	}

	if r.Update == nil {
		for _, s := range r.Schema {
			if s.ForceNew == false && (s.Computed == false || s.Optional == true) {
				r.Update = d.updateFn()
				break
			}
		}
	}
	return
}

func (d *description) Name() string {
	if d.name == "" {
		return "abiquo_" + d.media
	}
	return "abiquo_" + d.name
}
