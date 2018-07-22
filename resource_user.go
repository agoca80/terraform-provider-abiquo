package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var userSchema = map[string]*schema.Schema{
	"enterprise": endpoint("enterprise"),
	"active":     attribute(required, boolean),
	"email":      attribute(required, text),
	"name":       attribute(required, text),
	"nick":       attribute(required, text),
	"surname":    attribute(required, text),
	"password":   attribute(optional, text, computed),
	"scope":      attribute(optional, link("scope")),
	"role":       attribute(required, link("role")),
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

func userRead(d *resourceData, resource core.Resource) (err error) {
	user := resource.(*abiquo.User)
	d.Set("active", user.Active)
	d.Set("email", user.Email)
	d.Set("name", user.Name)
	d.Set("nick", user.Nick)
	d.Set("surname", user.Surname)
	return
}

var user = &description{
	Resource: &schema.Resource{Schema: userSchema},
	dto:      userNew,
	endpoint: endpointPath("enterprise", "/users"),
	media:    "user",
	read:     userRead,
}
