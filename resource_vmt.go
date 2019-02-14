package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var virtualmachinetemplateSchema = map[string]*schema.Schema{
	"repo":        endpoint("dcrepository"),
	"cpu":         attribute(required, natural),
	"name":        attribute(required, text),
	"description": attribute(optional, text),
	"ova":         attribute(required, text, forceNew),
	"ram":         attribute(required, natural),
	"icon":        attribute(optional, href),
}

func virtualmachinetemplateDTO(d *resourceData) core.Resource {
	return &abiquo.VirtualMachineTemplate{
		CPURequired: d.integer("cpu"),
		Name:        d.string("name"),
		Description: d.string("description"),
		IconURL:     d.string("icon"),
		RAMRequired: d.integer("ram"),
	}
}

func virtualmachinetemplateCreate(rd *schema.ResourceData, m interface{}) (err error) {
	d := newDataType(rd, "virtualmachinetemplate")
	endpoint := d.link("repo").SetType("datacenterrepository")
	resource, err := endpoint.Walk()
	if err != nil {
		return
	}

	dcrepo := resource.(*abiquo.DatacenterRepository)
	virtualmachinetemplate, err := dcrepo.UploadOVA(d.string("ova"))
	if err != nil {
		return
	}

	d.SetId(virtualmachinetemplate.URL())
	virtualmachinetemplate.Name = d.string("name")
	virtualmachinetemplate.IconURL = d.string("icon")
	virtualmachinetemplate.Description = d.string("description")
	virtualmachinetemplate.CPURequired = d.integer("cpu")
	virtualmachinetemplate.RAMRequired = d.integer("ram")
	err = virtualmachinetemplate.Update(virtualmachinetemplate)
	return
}

func virtualmachinetemplateRead(d *resourceData, resource core.Resource) (err error) {
	virtualmachinetemplate := resource.(*abiquo.VirtualMachineTemplate)
	d.Set("name", virtualmachinetemplate.Name)
	d.Set("icon", virtualmachinetemplate.IconURL)
	d.Set("description", virtualmachinetemplate.Description)
	d.Set("cpu", virtualmachinetemplate.CPURequired)
	d.Set("ram", virtualmachinetemplate.RAMRequired)
	return
}

var virtualmachinetemplate = &description{
	dto:   virtualmachinetemplateDTO,
	media: "virtualmachinetemplate",
	read:  virtualmachinetemplateRead,
	Resource: &schema.Resource{
		Schema: virtualmachinetemplateSchema,
		Create: virtualmachinetemplateCreate,
	},
}
