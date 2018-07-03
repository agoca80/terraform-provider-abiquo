package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var scopeDataSchema = map[string]*schema.Schema{
	"name": attribute(required, text),
}

func scopeDataRead(d *schema.ResourceData, meta interface{}) (err error) {
	scope := abiquo.Scopes(nil).Find(func(r core.Resource) bool {
		return r.(*abiquo.Scope).Name == d.Get("name").(string)
	})
	if scope == nil {
		return fmt.Errorf("scope %q was not found", d.Get("name"))
	}

	d.SetId(scope.URL())
	return
}
