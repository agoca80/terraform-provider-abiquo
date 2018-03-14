package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var userSchema = map[string]*schema.Schema{
	"active": &schema.Schema{
		Required: true,
		Type:     schema.TypeBool,
	},
	"email": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"name": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"nick": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"surname": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"password": &schema.Schema{
		Computed: true,
		Type:     schema.TypeString,
	},
	"enterprise": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
	"scope": &schema.Schema{
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
	"role": &schema.Schema{
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
}

var userResource = &schema.Resource{
	Schema:   userSchema,
	Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough},
	Read:     resourceRead(userNew, userRead, "user"),
	Create:   resourceCreate(userNew, nil, userRead, userEndpoint),
	Update:   resourceUpdate(userNew, nil, "user"),
	Exists:   resourceExists("user"),
	Delete:   resourceDelete,
}

func userNew(d *resourceData) core.Resource {
	return &abiquo.User{
		Active:   d.bool("active"),
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
