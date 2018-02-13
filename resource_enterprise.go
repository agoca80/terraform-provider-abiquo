package main

import (
	"github.com/abiquo/opal/abiquo"
	"github.com/abiquo/opal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var enterpriseResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name": Required().String(),
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
	},
	Delete: resourceDelete,
	Read:   enterpriseRead,
	Create: enterpriseCreate,
	Exists: resourceExists("enterprise"),
	Update: enterpriseUpdate,
}

func enterpriseDTO(rd *schema.ResourceData) *abiquo.Enterprise {
	d := newResourceData(rd, "")
	return &abiquo.Enterprise{
		Name:     d.string("name"),
		CPUHard:  d.int("cpuhard"),
		CPUSoft:  d.int("cpusoft"),
		HDHard:   d.int("hdhard"),
		HDSoft:   d.int("HDSoft"),
		IPHard:   d.int("iphard"),
		IPSoft:   d.int("ipsoft"),
		RAMHard:  d.int("ramhard"),
		RAMSoft:  d.int("ramsoft"),
		RepoSoft: d.int("reposoft"),
		RepoHard: d.int("repohard"),
		VolHard:  d.int("volhard"),
		VolSoft:  d.int("VolSoft"),
		VLANHard: d.int("vlanhard"),
		VLANSoft: d.int("vlansoft"),
	}
}

func enterpriseRead(rd *schema.ResourceData, m interface{}) (err error) {
	e := new(abiquo.Enterprise)
	d := newResourceData(rd, "enterprise")
	if err = core.Read(core.NewLinkType(d.Id(), "enterprise"), e); err == nil {
		d.Set("name", e.Name)
		d.Set("cpuhard", e.CPUHard)
		d.Set("cpusoft", e.CPUSoft)
		d.Set("hdhard", e.HDHard)
		d.Set("hdsoft", e.HDSoft)
		d.Set("ipsoft", e.IPSoft)
		d.Set("iphard", e.IPHard)
		d.Set("ramsoft", e.RAMSoft)
		d.Set("ramhard", e.RAMHard)
		d.Set("reposoft", e.RepoSoft)
		d.Set("repohard", e.RepoHard)
		d.Set("volhard", e.VolHard)
		d.Set("volsoft", e.VolSoft)
		d.Set("vlanhard", e.VLANHard)
		d.Set("vlansoft", e.VLANSoft)
	}
	return
}

func enterpriseCreate(rd *schema.ResourceData, m interface{}) (err error) {
	e := enterpriseDTO(rd)
	if err = e.Create(); err == nil {
		rd.SetId(e.URL())
	}
	return
}

func enterpriseUpdate(rd *schema.ResourceData, m interface{}) (err error) {
	e := enterpriseDTO(rd)
	return core.Update(e, e)
}
