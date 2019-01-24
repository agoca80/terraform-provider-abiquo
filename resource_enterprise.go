package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var enterpriseSchema = map[string]*schema.Schema{
	"name":            attribute(required, text),
	"properties":      attribute(optional, hash(text)),
	"pricingtemplate": attribute(optional, href),
	"cpuhard":         attribute(optional, natural),
	"cpusoft":         attribute(optional, natural),
	"hdhard":          attribute(optional, natural),
	"hdsoft":          attribute(optional, natural),
	"iphard":          attribute(optional, natural),
	"ipsoft":          attribute(optional, natural),
	"ramhard":         attribute(optional, natural),
	"ramsoft":         attribute(optional, natural),
	"repohard":        attribute(optional, natural),
	"reposoft":        attribute(optional, natural),
	"vlanhard":        attribute(optional, natural),
	"volsoft":         attribute(optional, natural),
	"volhard":         attribute(optional, natural),
	"vlansoft":        attribute(optional, natural),
}

func enterpriseDTO(d *resourceData) core.Resource {
	return &abiquo.Enterprise{
		Name:     d.string("name"),
		CPUHard:  d.integer("cpuhard"),
		CPUSoft:  d.integer("cpusoft"),
		HDHard:   d.integer("hdhard"),
		HDSoft:   d.integer("HDSoft"),
		IPHard:   d.integer("iphard"),
		IPSoft:   d.integer("ipsoft"),
		RAMHard:  d.integer("ramhard"),
		RAMSoft:  d.integer("ramsoft"),
		RepoSoft: d.integer("reposoft"),
		RepoHard: d.integer("repohard"),
		VolHard:  d.integer("volhard"),
		VolSoft:  d.integer("VolSoft"),
		VLANHard: d.integer("vlanhard"),
		VLANSoft: d.integer("vlansoft"),
		DTO: core.NewDTO(
			d.link("pricingtemplate"),
		),
	}
}

func enterpriseRead(d *resourceData, resource core.Resource) (err error) {
	e := resource.(*abiquo.Enterprise)
	props, err := e.Rel("properties").Walk()
	if err != nil {
		return
	}
	properties := props.(*abiquo.EnterpriseProperties)
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
	d.Set("pricingtemplate", e.Rel("pricingtemplate").URL())
	return
}

func enterpriseUpdate(d *resourceData, enterprise core.Resource) (err error) {
	if d.HasChange("properties") {
		err = core.Update(enterprise.Rel("properties"), enterpriseProperties(d))
	}
	return
}

func enterpriseProperties(d *resourceData) *abiquo.EnterpriseProperties {
	properties := new(abiquo.EnterpriseProperties)
	properties.Properties = make(map[string]string)
	for k, v := range d.dict("properties") {
		properties.Properties[k] = v.(string)
	}
	return properties
}

var enterprise = &description{
	Resource: &schema.Resource{Schema: enterpriseSchema},
	dto:      enterpriseDTO,
	endpoint: endpointConst("admin/enterprises"),
	media:    "enterprise",
	create:   enterpriseUpdate,
	update:   enterpriseUpdate,
	read:     enterpriseRead,
}
