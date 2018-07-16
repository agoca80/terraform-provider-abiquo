package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var vmtSchema = map[string]*schema.Schema{
	"repo":        endpoint("dcrepository"),
	"cpu":         attribute(required, natural),
	"name":        attribute(required, text),
	"description": attribute(optional, text),
	"file":        attribute(required, text, forceNew),
	"ram":         attribute(required, natural),
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
	endpoint := d.link("repo").SetType("datacenterrepository")
	resource := endpoint.Walk()
	if resource == nil {
		return fmt.Errorf("repository %q does not exist", d.string("repo"))
	}

	dcrepo := resource.(*abiquo.DatacenterRepository)
	vmt, err := dcrepo.Upload(d.string("file"))
	if err != nil {
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
