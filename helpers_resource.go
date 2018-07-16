package main

import (
	"fmt"

	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

type (
	factory        func(*resourceData) core.Resource
	resourceMethod func(*resourceData, core.Resource) error
)

func resourceDelete(d *schema.ResourceData, m interface{}) (err error) {
	return core.Remove(newResourceData(d, ""))
}

func resourceExists(media string) schema.ExistsFunc {
	return func(d *schema.ResourceData, m interface{}) (ok bool, err error) {
		return core.Read(newResourceData(d, media), nil) == nil, nil
	}
}

func resourceCreate(factory factory, create resourceMethod, read resourceMethod, endpoint func(*resourceData) *core.Link) schema.CreateFunc {
	return func(rd *schema.ResourceData, m interface{}) (err error) {
		d := newResourceData(rd, "")

		resource := factory(d)
		if resource == nil {
			return fmt.Errorf("resourceCreate: resource could not be created")
		}

		if err = core.Create(endpoint(d), resource); err != nil {
			return
		}
		d.SetId(resource.URL())
		if create != nil {
			err = create(d, resource)
		}

		if err == nil && read != nil {
			err = read(d, resource)
		}

		return
	}
}

func resourceUpdate(factory factory, update resourceMethod, media string) schema.UpdateFunc {
	return func(rd *schema.ResourceData, m interface{}) (err error) {
		d := newResourceData(rd, media)
		resource := factory(d)

		if err = core.Update(d, resource); err == nil && update != nil {
			err = update(d, resource)
		}

		return
	}
}

func resourceRead(factory factory, read resourceMethod, media string) schema.ReadFunc {
	return func(rd *schema.ResourceData, m interface{}) (err error) {
		d := newResourceData(rd, media)
		resource := factory(d)
		if resource == nil {
			err = fmt.Errorf("resourceRead: could not create %v resource", media)
		} else if err = core.Read(d, resource); err == nil {
			err = read(d, resource)
		}
		return
	}
}
