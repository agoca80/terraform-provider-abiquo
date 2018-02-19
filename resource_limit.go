package main

import (
	"github.com/abiquo/opal/abiquo"
	"github.com/abiquo/opal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var limitResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		// Soft limits
		"cpusoft":  Optional().Number(),
		"hdsoft":   Optional().Number(),
		"ipsoft":   Optional().Number(),
		"ramsoft":  Optional().Number(),
		"reposoft": Optional().Number(),
		"volsoft":  Optional().Number(),
		"vlansoft": Optional().Number(),
		// Hard limits
		"cpuhard":  Optional().Number(),
		"hdhard":   Optional().Number(),
		"iphard":   Optional().Number(),
		"ramhard":  Optional().Number(),
		"repohard": Optional().Number(),
		"vlanhard": Optional().Number(),
		"volhard":  Optional().Number(),
		// Links
		"location":   Required().Renew().Link(),
		"enterprise": Required().Renew().Link(),
		"hwprofiles": Optional().Links(),
	},
	Exists: resourceExists("limit"),
	Read:   resourceRead(limitNew, limitRead, "limit"),
	Update: resourceUpdate(limitNew, "limit"),
	Create: resourceCreate(limitNew, nil, limitRead, limitEndpoint),
	Delete: resourceDelete,
}

func limitNew(d *resourceData) core.Resource {
	return &abiquo.Limit{
		// Soft limits
		CPUSoft:  d.int("cpusoft"),
		HDSoft:   d.int("hdsoft"),
		IPSoft:   d.int("ipsoft"),
		RAMSoft:  d.int("ramsoft"),
		RepoSoft: d.int("reposoft"),
		VolSoft:  d.int("VolSoft"),
		VLANSoft: d.int("vlansoft"),
		// Hard limits
		CPUHard:  d.int("cpuhard"),
		HDHard:   d.int("hdhard"),
		IPHard:   d.int("iphard"),
		RAMHard:  d.int("ramhard"),
		RepoHard: d.int("repohard"),
		VolHard:  d.int("volhard"),
		VLANHard: d.int("vlanhard"),
		// Links
		DTO: core.NewDTO(
			d.link("location"),
		),
	}
}

func limitEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("enterprise")+"/limits", "limit")
}

func limitRead(d *resourceData, resource core.Resource) (err error) {
	limit := resource.(*abiquo.Limit)
	// Soft limits
	d.Set("cpusoft", limit.CPUSoft)
	d.Set("hdsoft", limit.HDSoft)
	d.Set("ipsoft", limit.IPSoft)
	d.Set("ramsoft", limit.RAMSoft)
	d.Set("reposoft", limit.RepoSoft)
	d.Set("volsoft", limit.VolSoft)
	d.Set("vlansoft", limit.VLANSoft)
	// Hard limits
	d.Set("cpuhard", limit.CPUHard)
	d.Set("hdhard", limit.HDHard)
	d.Set("iphard", limit.IPHard)
	d.Set("ramhard", limit.RAMHard)
	d.Set("repohard", limit.RepoHard)
	d.Set("volhard", limit.VolHard)
	d.Set("vlanhard", limit.VLANHard)
	return
}
