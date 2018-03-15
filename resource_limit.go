package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

var limitSchema = map[string]*schema.Schema{
	// Soft limits
	"cpusoft": &schema.Schema{
		Optional:     true,
		Type:         schema.TypeInt,
		ValidateFunc: validation.IntAtLeast(0),
	},
	"hdsoft": &schema.Schema{
		Optional:     true,
		Type:         schema.TypeInt,
		ValidateFunc: validation.IntAtLeast(0),
	},
	"ipsoft": &schema.Schema{
		Optional:     true,
		Type:         schema.TypeInt,
		ValidateFunc: validation.IntAtLeast(0),
	},
	"ramsoft": &schema.Schema{
		Optional:     true,
		Type:         schema.TypeInt,
		ValidateFunc: validation.IntAtLeast(0),
	},
	"reposoft": &schema.Schema{
		Optional:     true,
		Type:         schema.TypeInt,
		ValidateFunc: validation.IntAtLeast(0),
	},
	"volsoft": &schema.Schema{
		Optional:     true,
		Type:         schema.TypeInt,
		ValidateFunc: validation.IntAtLeast(0),
	},
	"vlansoft": &schema.Schema{
		Optional:     true,
		Type:         schema.TypeInt,
		ValidateFunc: validation.IntAtLeast(0),
	},
	// Hard limits
	"cpuhard": &schema.Schema{
		Optional:     true,
		Type:         schema.TypeInt,
		ValidateFunc: validation.IntAtLeast(0),
	},
	"hdhard": &schema.Schema{
		Optional:     true,
		Type:         schema.TypeInt,
		ValidateFunc: validation.IntAtLeast(0),
	},
	"iphard": &schema.Schema{
		Optional:     true,
		Type:         schema.TypeInt,
		ValidateFunc: validation.IntAtLeast(0),
	},
	"ramhard": &schema.Schema{
		Optional:     true,
		Type:         schema.TypeInt,
		ValidateFunc: validation.IntAtLeast(0),
	},
	"repohard": &schema.Schema{
		Optional:     true,
		Type:         schema.TypeInt,
		ValidateFunc: validation.IntAtLeast(0),
	},
	"vlanhard": &schema.Schema{
		Optional:     true,
		Type:         schema.TypeInt,
		ValidateFunc: validation.IntAtLeast(0),
	},
	"volhard": &schema.Schema{
		Optional:     true,
		Type:         schema.TypeInt,
		ValidateFunc: validation.IntAtLeast(0),
	},
	// Links
	"location": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
	"enterprise": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
	"hwprofiles": &schema.Schema{
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validateURL,
		},
		Optional: true,
		Set:      schema.HashString,
		Type:     schema.TypeSet,
	},
	"backups": &schema.Schema{
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validateURL,
		},
		Optional: true,
		Set:      schema.HashString,
		Type:     schema.TypeSet,
	},
	"dstiers": &schema.Schema{
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validateURL,
		},
		Optional: true,
		Set:      schema.HashString,
		Type:     schema.TypeSet,
	},
}

func limitNew(d *resourceData) core.Resource {
	limit := &abiquo.Limit{
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

	// Backups
	backups := d.set("backups")
	if backups != nil && backups.Len() > 0 {
		for _, entry := range backups.List() {
			href := entry.(string)
			limit.Add(core.NewLinkType(href, "backuppolicy").SetRel("backuppolicy"))
		}
	}

	// HWprofiles
	hwprofiles := d.set("hwprofiles")
	if hwprofiles != nil && hwprofiles.Len() > 0 {
		for _, entry := range hwprofiles.List() {
			href := entry.(string)
			limit.Add(core.NewLinkType(href, "hardwareprofile").SetRel("hardwareprofile"))
		}
	}

	// DSTiers
	dstiers := d.set("dstiers")
	if dstiers != nil && dstiers.Len() > 0 {
		for _, entry := range dstiers.List() {
			href := entry.(string)
			limit.Add(core.NewLinkType(href, "datastoretier").SetRel("datastoretier"))
		}
	}

	return limit
}

func limitEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("enterprise")+"/limits", "limit")
}

func limitRead(d *resourceData, resource core.Resource) (err error) {
	limit := resource.(*abiquo.Limit)

	backups := mapHrefs(limit.Links.Filter(func(l *core.Link) bool {
		return l.IsMedia("backuppolicy")
	}))

	hwprofiles := mapHrefs(limit.Links.Filter(func(l *core.Link) bool {
		return l.IsMedia("hwprofile")
	}))

	dstiers := mapHrefs(limit.Links.Filter(func(l *core.Link) bool {
		return l.IsMedia("datastoretier")
	}))

	d.Set("backups", backups)
	d.Set("hwprofiles", hwprofiles)
	d.Set("dstiers", dstiers)
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
