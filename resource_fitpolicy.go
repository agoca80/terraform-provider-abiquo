package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var fitPolicySchema = map[string]*schema.Schema{
	"policy": attribute(required, forceNew, label([]string{"PROGRESSIVE", "PERFORMANCE"})),
	"target": attribute(required, forceNew, href),
}

func fitPolicyDTO(d *resourceData) core.Resource {
	return &abiquo.FitPolicy{
		FitPolicy: d.string("policy"),
		DTO: core.NewDTO(
			d.linkTypeRel("target", "datacenter", "datacenter"),
		),
	}
}

func fitPolicyEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType("admin/rules/fitsPolicy", "fitpolicyrule")
}

func fitPolicyRead(d *resourceData, resource core.Resource) (err error) {
	policy := resource.(*abiquo.FitPolicy)
	d.Set("policy", policy.FitPolicy)
	d.Set("target", policy.Rel("datacenter").URL())
	return
}
