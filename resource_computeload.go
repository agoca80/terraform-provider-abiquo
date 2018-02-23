package main

import (
	"strings"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var computeLoadSchema = map[string]*schema.Schema{
	"aggregated": Optional().Bool(),
	"cpuload":    Required().Number(),
	"ramload":    Required().Number(),
	"target":     Optional().Link(),
}

var computeLoadResource = &schema.Resource{
	Schema: computeLoadSchema,
	Delete: resourceDelete,
	Exists: resourceExists("machineloadrule"),
	Create: resourceCreate(computeLoadDTO, nil, computeLoadRead, computeLoadEndpoint),
	Update: resourceUpdate(computeLoadDTO, "machineloadrule"),
	Read:   resourceRead(computeLoadDTO, computeLoadRead, "machineloadrule"),
}

func computeLoadDTO(d *resourceData) core.Resource {
	machineLoadRule := &abiquo.MachineLoadRule{
		Aggregated:        d.bool("aggregated"),
		CPULoadPercentage: d.int("cpuload"),
		RAMLoadPercentage: d.int("ramload"),
	}

	if h, ok := d.GetOk("target"); ok {
		var media string
		var href = h.(string)
		switch {
		case strings.Contains(href, "cluster"):
			media = "cluster"
		case strings.Contains(href, "machine"):
			media = "machine"
		case strings.Contains(href, "rack"):
			media = "rack"
		case strings.Contains(href, "datacenter"):
			media = "datacenter"
		default:
			return nil
		}
		machineLoadRule.Add(core.NewLinkType(href, media).SetRel(media))
	}

	return machineLoadRule
}

func computeLoadEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType("admin/rules/machineLoadLevel", "machineloadrule")
}

func computeLoadRead(d *resourceData, resource core.Resource) (err error) {
	rule := resource.(*abiquo.MachineLoadRule)
	d.Set("aggregated", rule.Aggregated)
	d.Set("cpuload", rule.CPULoadPercentage)
	d.Set("ramload", rule.RAMLoadPercentage)
	for _, media := range []string{"cluster", "machine", "rack", "datacenter"} {
		if rel := rule.Rel(media); rel != nil {
			d.Set("target", rel.URL())
		}
	}
	return
}
