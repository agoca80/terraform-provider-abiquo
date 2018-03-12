package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var alertResource = &schema.Resource{
	Schema: alertSchema,
	Delete: resourceDelete,
	Exists: resourceExists("alert"),
	Update: resourceUpdate(alertNew, nil, "alert"),
	Create: resourceCreate(alertNew, nil, alertRead, alertEndpoint),
	Read:   resourceRead(alertNew, alertRead, "alert"),
}

var alertSchema = map[string]*schema.Schema{
	"virtualappliance": Required().Renew().Link(),
	"name":             Required().String(),
	"description":      Required().String(),
	"alarms":           Required().Links(),
}

func alertNew(d *resourceData) core.Resource {
	slice := d.slice("alarms")
	alarms := make([]*core.Link, len(slice))
	for i, a := range slice {
		alarms[i] = core.NewLinkType(a.(string), "alarm")
	}

	return &abiquo.Alert{
		Name:        d.string("name"),
		Description: d.string("description"),
		Alarms:      alarms,
	}
}

func alertEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("virtualappliance")+"/alerts", "alert")
}

func alertRead(d *resourceData, resource core.Resource) (err error) {
	alert := resource.(*abiquo.Alert)
	// PENDING alarms are unordered, so, set comparation must be performed
	// to detect changed
	// d.Set("alarms", ...)
	d.Set("name", alert.Name)
	d.Set("description", alert.Description)
	return
}
