package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var licenseSchema = map[string]*schema.Schema{
	"code":       attribute(required, text, forceNew),
	"expiration": attribute(computed, text),
	"numcores":   attribute(computed, integer),
	"sgenabled":  attribute(computed, boolean),
}

func licenseNew(d *resourceData) core.Resource {
	return &abiquo.License{
		Code: d.string("code"),
	}
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

var license = &description{
	dto:      licenseNew,
	endpoint: endpointConst("config/licenses"),
	media:    "license",
	read:     licenseRead,
	Resource: &schema.Resource{
		Schema: licenseSchema,
	},
}
