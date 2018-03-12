package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var enterpriseSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	// Soft limits
	"cpusoft": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	"hdsoft": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	"ipsoft": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	"ramsoft": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	"reposoft": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	"volsoft": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	"vlansoft": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	// Hard limits
	"cpuhard": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	"hdhard": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	"iphard": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	"ramhard": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	"repohard": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	"vlanhard": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
	"volhard": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
}

var enterpriseResource = &schema.Resource{
	Schema: enterpriseSchema,
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
