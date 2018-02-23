package main

import (
	"fmt"
	"net/url"

	"github.com/abiquo/ojal/abiquo"
	"github.com/hashicorp/terraform/helper/schema"
)

var enterpriseDataSource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name": Required().String(),
	},
	Read: enterpriseDataRead,
}

func enterpriseDataRead(d *schema.ResourceData, p interface{}) (err error) {
	query := url.Values{"has": {d.Get("name").(string)}}
	enterprise := abiquo.Enterprises(query).First()
	if enterprise == nil {
		return fmt.Errorf("enterprise %q was not found", d.Get("name"))
	}

	d.SetId(enterprise.URL())
	return
}
