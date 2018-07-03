package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var templateDataSchema = map[string]*schema.Schema{
	"name": attribute(required, text),
}

func templateRead(d *schema.ResourceData, meta interface{}) (err error) {
	enterprise := meta.(*provider).Enterprise()
	dcrepos := enterprise.Rel("datacenterrepositories").Collection(nil)
	for _, repo := range dcrepos.List() {
		vmtemplates := repo.Rel("virtualmachinetemplates").Collection(nil)
		template := vmtemplates.Find(func(r core.Resource) bool {
			t := r.(*abiquo.VirtualMachineTemplate)
			return t.Name == d.Get("name").(string) && t.State != "UNAVAILABLE"
		})
		if template != nil {
			d.SetId(template.URL())
			return
		}
	}

	return fmt.Errorf("template %q was not found", d.Get("name"))
}
