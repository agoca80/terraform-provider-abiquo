package main

import (
	"github.com/abiquo/opal/abiquo"
	"github.com/abiquo/opal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var fitPolicySchema = map[string]*schema.Schema{
	"policy": Required().Renew().ValidateString([]string{"PROGRESSIVE", "PERFORMANCE"}),
	"target": Required().Renew().Link(),
}

var fitPolicyResource = &schema.Resource{
	Schema: fitPolicySchema,
	Delete: resourceDelete,
	Exists: resourceExists("fitpolicyrule"),
	Create: resourceCreate(fitPolicyDTO, nil, fitPolicyRead, fitPolicyEndpoint),
	Read:   resourceRead(fitPolicyDTO, fitPolicyRead, "fitpolicyrule"),
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
