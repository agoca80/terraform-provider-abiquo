package main

import (
	"fmt"

	"github.com/abiquo/opal/abiquo"
	"github.com/abiquo/opal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var scopeDataSource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name": Required().String(),
	},
	Read: scopeDataRead,
}

func scopeDataRead(d *schema.ResourceData, meta interface{}) (err error) {
	finder := func(r core.Resource) bool {
		return r.(*abiquo.Scope).Name == d.Get("name").(string)
	}
	if scope := abiquo.Scopes(nil).Find(finder); scope != nil {
		d.SetId(scope.URL())
		return
	}

	return fmt.Errorf("vdc %q was not found", d.Get("name"))
}
