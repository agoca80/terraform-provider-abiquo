package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var licenseResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"code": &schema.Schema{
			Required: true,
			ForceNew: true,
			Type:     schema.TypeString,
		},
		"expiration": &schema.Schema{
			Computed: true,
			Type:     schema.TypeString,
		},
		"numcores": &schema.Schema{
			Computed: true,
			Type:     schema.TypeInt,
		},
		"sgenabled": &schema.Schema{
			Computed: true,
			Type:     schema.TypeBool,
		},
	},
	Delete: resourceDelete,
	Exists: resourceExists("license"),
	Create: resourceCreate(licenseNew, nil, licenseRead, licenseEndpoint),
	Read:   resourceRead(licenseNew, licenseRead, "license"),
}

func licenseNew(d *resourceData) core.Resource {
	return &abiquo.License{
		Code: d.string("code"),
	}
}

func licenseEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType("config/licenses", "license")
}

func licenseRead(d *resourceData, resource core.Resource) (err error) {
	license := resource.(*abiquo.License)
	d.Set("id", license.ID)
	d.Set("code", license.Code)
	d.Set("expiration", license.Expiration)
	d.Set("numcores", license.NumCores)
	d.Set("sgenabled", license.ScalingGroupsEnabled)
	return
}
