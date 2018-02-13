package main

import (
	"errors"
	"fmt"

	"github.com/abiquo/opal/abiquo"
	"github.com/abiquo/opal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var templateDataSource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name": Required().String(),
	},
	Read: templateRead,
}

func templateRead(d *schema.ResourceData, meta interface{}) (err error) {
	enterprise := meta.(*provider).Enterprise()
	if enterprise == nil {
		return errors.New("The user enterprise was not found")
	}

	for _, dcrepo := range enterprise.DatacenterRepositories(nil).List() {
		repo := dcrepo.(*abiquo.DatacenterRepository)
		finder := func(r core.Resource) bool {
			return r.(*abiquo.VirtualMachineTemplate).Name == d.Get("name").(string)
		}
		if template := repo.VirtualMachineTemplates(nil).Find(finder); template != nil {
			d.SetId(template.URL())
			return
		}
	}

	return fmt.Errorf("template %q was not found", d.Get("name"))
}
