package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var vmtSchema = map[string]*schema.Schema{
	"cpu":         attribute(required, natural),
	"name":        attribute(required, text),
	"description": attribute(optional, text),
	"file":        attribute(required, text, forceNew),
	"ram":         attribute(required, natural),
	"repo":        attribute(required, href, forceNew),
	"icon":        attribute(optional, href),
}

func vmtNew(d *resourceData) core.Resource {
	return &abiquo.VirtualMachineTemplate{
		CPURequired: d.integer("cpu"),
		Name:        d.string("name"),
		Description: d.string("description"),
		IconURL:     d.string("icon"),
		RAMRequired: d.integer("ram"),
	}
}

func vmtCreate(rd *schema.ResourceData, m interface{}) (err error) {
	d := newResourceData(rd, "virtualmachinetemplate")
	file := d.string("file")
	endpoint := d.link("repo").SetType("datacenterrepository")
	repository := new(abiquo.DatacenterRepository)
	if err = core.Read(endpoint, repository); err != nil {
		return
	}

	var vmt *abiquo.VirtualMachineTemplate
	if vmt, err = repository.Upload(file); err != nil {
		return
	}

	d.SetId(vmt.URL())
	vmt.Name = d.string("name")
	vmt.IconURL = d.string("icon")
	vmt.Description = d.string("description")
	vmt.CPURequired = d.integer("cpu")
	vmt.RAMRequired = d.integer("ram")
	err = core.Update(vmt, vmt)

	return
}

func vmtRead(d *resourceData, resource core.Resource) (err error) {
	vmt := resource.(*abiquo.VirtualMachineTemplate)
	d.Set("name", vmt.Name)
	d.Set("icon", vmt.IconURL)
	d.Set("description", vmt.Description)
	d.Set("cpu", vmt.CPURequired)
	d.Set("ram", vmt.RAMRequired)
	return
}

func vmtUpdate(rd *schema.ResourceData, meta interface{}) (err error) {
	d := newResourceData(rd, "virtualmachinetemplate")
	vmt := d.Walk().(*abiquo.VirtualMachineTemplate)
	vmt.Name = d.string("name")
	vmt.Description = d.string("description")
	vmt.IconURL = d.string("icon")
	vmt.CPURequired = d.integer("cpu")
	vmt.RAMRequired = d.integer("ram")
	return core.Update(d, vmt)
}
