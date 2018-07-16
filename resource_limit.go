package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var limitSchema = map[string]*schema.Schema{
	"enterprise": endpoint("enterprise"),
	"backups":    attribute(optional, set(link("backuppolicy_dc"))),
	"cpuhard":    attribute(optional, natural),
	"cpusoft":    attribute(optional, natural),
	"dstiers":    attribute(optional, set(link("datastoretier_dc"))),
	"hdhard":     attribute(optional, natural),
	"hdsoft":     attribute(optional, natural),
	"hwprofiles": attribute(optional, set(link("hardwareprofile_dc"))),
	"iphard":     attribute(optional, natural),
	"ipsoft":     attribute(optional, natural),
	"location":   attribute(required, forceNew, link("datacenter")),
	"ramhard":    attribute(optional, natural),
	"ramsoft":    attribute(optional, natural),
	"repohard":   attribute(optional, natural),
	"reposoft":   attribute(optional, natural),
	"volhard":    attribute(optional, natural),
	"volsoft":    attribute(optional, natural),
	"vlanhard":   attribute(optional, natural),
	"vlansoft":   attribute(optional, natural),
}

func limitNew(d *resourceData) core.Resource {
	limit := &abiquo.Limit{
		// Soft limits
		CPUSoft:  d.integer("cpusoft"),
		HDSoft:   d.integer("hdsoft"),
		IPSoft:   d.integer("ipsoft"),
		RAMSoft:  d.integer("ramsoft"),
		RepoSoft: d.integer("reposoft"),
		VolSoft:  d.integer("VolSoft"),
		VLANSoft: d.integer("vlansoft"),
		// Hard limits
		CPUHard:  d.integer("cpuhard"),
		HDHard:   d.integer("hdhard"),
		IPHard:   d.integer("iphard"),
		RAMHard:  d.integer("ramhard"),
		RepoHard: d.integer("repohard"),
		VolHard:  d.integer("volhard"),
		VLANHard: d.integer("vlanhard"),
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
		limit.EnableHPs = true
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
		return l.IsMedia("hardwareprofile")
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
