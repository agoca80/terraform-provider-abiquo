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
	"properties": &schema.Schema{
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional: true,
		Type:     schema.TypeMap,
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

func enterpriseDTO(d *resourceData) core.Resource {
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

func enterpriseEndpoint(_ *resourceData) *core.Link {
	return core.NewLinkType("admin/enterprises", "enterprise")
}

func enterpriseCreate(d *resourceData, enterprise core.Resource) (err error) {
	properties := enterpriseProperties(d)
	if len(properties.Properties) > 0 {
		err = core.Update(enterprise.Rel("properties"), properties)
	}
	return
}

func enterpriseRead(d *resourceData, resource core.Resource) (err error) {
	e := resource.(*abiquo.Enterprise)
	properties := e.Rel("properties").Walk().(*abiquo.EnterpriseProperties)
	d.Set("properties", properties.Properties)
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
	return
}

func enterpriseUpdate(d *resourceData, enterprise core.Resource) (err error) {
	if !d.HasChange("properties") {
		return
	}

	return core.Update(enterprise.Rel("properties"), enterpriseProperties(d))
}

func enterpriseProperties(d *resourceData) *abiquo.EnterpriseProperties {
	properties := new(abiquo.EnterpriseProperties)
	properties.Properties = make(map[string]string)
	for k, v := range d.dict("properties") {
		properties.Properties[k] = v.(string)
	}
	return properties
}
