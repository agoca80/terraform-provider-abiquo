package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var userSchema = map[string]*schema.Schema{
	"active":     attribute(required, boolean),
	"email":      attribute(required, text),
	"name":       attribute(required, text),
	"nick":       attribute(required, text),
	"surname":    attribute(required, text),
	"password":   attribute(computed, text),
	"enterprise": attribute(required, enterprise, forceNew),
	"scope":      attribute(optional, href),
	"role":       attribute(required, href),
}

func userNew(d *resourceData) core.Resource {
	return &abiquo.User{
		Active:   d.boolean("active"),
		Email:    d.string("email"),
		Name:     d.string("name"),
		Nick:     d.string("nick"),
		Password: "12qwaszx",
		Surname:  d.string("surname"),
		DTO: core.NewDTO(
			d.link("enterprise"),
			d.link("scope"),
			d.link("role"),
		),
	}
}

func userEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("enterprise")+"/users", "user")
}

func userRead(d *resourceData, resource core.Resource) (err error) {
	user := resource.(*abiquo.User)
	d.Set("active", user.Active)
	d.Set("email", user.Email)
	d.Set("name", user.Name)
	d.Set("nick", user.Nick)
	d.Set("surname", user.Surname)
	return
}
