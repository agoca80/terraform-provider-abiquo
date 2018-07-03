package main

import (
	"fmt"

	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var repoDataSchema = map[string]*schema.Schema{
	"datacenter": attribute(required, text),
}

func dataRepoRead(d *schema.ResourceData, p interface{}) (err error) {
	enterprise := p.(*provider).Enterprise()
	repos := enterprise.Rel("datacenterrepositories").Collection(nil)
	repo := repos.Find(func(r core.Resource) bool {
		return r.Rel("datacenter").Title == d.Get("datacenter").(string)
	})
	if repo == nil {
		return fmt.Errorf("datacenter repository for datacenter %q was not found", d.Get("datacenter"))
	}
	d.SetId(repo.URL())
	return
}
