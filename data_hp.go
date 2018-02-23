package main

import (
	"fmt"
	"net/url"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"

	"github.com/hashicorp/terraform/helper/schema"
)

var hpDataSource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name":     Required().String(),
		"location": Required().Link(),
	},
	Read: hpDataRead,
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
