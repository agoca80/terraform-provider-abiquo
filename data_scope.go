package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var scopeDataSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
}

var scopeDataSource = &schema.Resource{
	Schema: scopeDataSchema,
	Read:   scopeDataRead,
}

func scopeDataRead(d *schema.ResourceData, meta interface{}) (err error) {
	finder := func(r core.Resource) bool {
		return r.(*abiquo.Scope).Name == d.Get("name").(string)
	}
	if scope := abiquo.Scopes(nil).Find(finder); scope != nil {
		d.SetId(scope.URL())
		return
	}

	return fmt.Errorf("scope %q was not found", d.Get("name"))
}
