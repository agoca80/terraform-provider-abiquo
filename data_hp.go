package main

import (
	"fmt"
	"net/url"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"

	"github.com/hashicorp/terraform/helper/schema"
)

var hpDataSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"location": &schema.Schema{
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
}

var hpDataSource = &schema.Resource{
	Schema: hpDataSchema,
	Read:   hpDataRead,
}

func hpDataRead(d *schema.ResourceData, meta interface{}) (err error) {
	locationName := d.Get("location").(string)
	finder := func(r core.Resource) bool {
		return r.(*abiquo.Datacenter).Name == locationName
	}
	l := abiquo.PrivateLocations(nil).Find(finder)
	if l == nil {
		return fmt.Errorf("location %q does not exist", locationName)
	}

	profileName := d.Get("name").(string)
	query := url.Values{"active": {"true"}, "has": {profileName}}
	hp := l.Rel("hardwareprofiles").Collection(query).First()
	if hp == nil {
		return fmt.Errorf("hwprofile %q does not exist in %q", profileName, locationName)
	}
	d.SetId(hp.URL())

	return
}
