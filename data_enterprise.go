package main

import (
	"fmt"
	"net/url"

	"github.com/abiquo/ojal/abiquo"
	"github.com/hashicorp/terraform/helper/schema"
)

var enterpriseDataSchema = map[string]*schema.Schema{
	"name": attribute(required, text),
}

func enterpriseFind(d *resourceData) (err error) {
	query := url.Values{"has": {d.string("name")}}
	enterprise := abiquo.Enterprises(query).First()
	if enterprise == nil {
		return fmt.Errorf("enterprise %q was not found", d.Get("name"))
	}

	d.SetId(enterprise.URL())
	return
}
