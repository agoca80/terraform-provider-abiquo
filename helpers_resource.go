package main

import (
	"fmt"

	"github.com/abiquo/opal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

type factory func(*resourceData) core.Resource
type endpoint func(*resourceData) *core.Link
type readResource func(*resourceData, core.Resource) error

func resourceDelete(d *schema.ResourceData, m interface{}) (err error) {
	return core.Remove(newResourceData(d, ""))
}

func resourceExists(media string) schema.ExistsFunc {
	return func(d *schema.ResourceData, m interface{}) (ok bool, err error) {
		return core.Read(newResourceData(d, media), nil) == nil, nil
	}
}

func resourceCreate(factory factory, endpoint endpoint) schema.CreateFunc {
	return func(rd *schema.ResourceData, m interface{}) (err error) {
		d := newResourceData(rd, "")

		resource := factory(d)
		if resource == nil {
			return fmt.Errorf("resourceCreate: resource could not be created")
		}

		if err = core.Create(endpoint(d), resource); err == nil {
			d.SetId(resource.URL())
		}
		return
	}
}

func resourceUpdate(factory factory, media string) schema.UpdateFunc {
	return func(rd *schema.ResourceData, m interface{}) (err error) {
		d := newResourceData(rd, media)
		resource := factory(d)
		return core.Update(d, resource)
	}
}

func resourceRead(factory factory, readResource readResource, media string) schema.ReadFunc {
	return func(rd *schema.ResourceData, m interface{}) (err error) {
		d := newResourceData(rd, media)
		resource := factory(d)
		if resource == nil {
			err = fmt.Errorf("resourceRead: could not create %v resource", media)
		} else if err = core.Read(d, resource); err == nil {
			err = readResource(d, resource)
		}
		return
	}
}
