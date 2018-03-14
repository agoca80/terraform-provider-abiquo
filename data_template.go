package main

import (
	"errors"
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var templateDataSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
}

var templateDataSource = &schema.Resource{
	Schema: templateDataSchema,
	Read:   templateRead,
}

func templateRead(d *schema.ResourceData, meta interface{}) (err error) {
	enterprise := meta.(*provider).Enterprise()
	if enterprise == nil {
		return errors.New("The user enterprise was not found")
	}

	dcrepos := enterprise.Rel("datacenterrepositories").Collection(nil)
	for _, dcrepo := range dcrepos.List() {
		repo := dcrepo.(*abiquo.DatacenterRepository)
		finder := func(r core.Resource) bool {
			return r.(*abiquo.VirtualMachineTemplate).Name == d.Get("name").(string)
		}
		vmtemplates := repo.Rel("virtualmachinetemplates").Collection(nil)
		if template := vmtemplates.Find(finder); template != nil {
			d.SetId(template.URL())
			return
		}
	}

	return fmt.Errorf("template %q was not found", d.Get("name"))
}
