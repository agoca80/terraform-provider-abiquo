package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var vmSchema = map[string]*schema.Schema{
	"cpu":                    attribute(optional, forceNew, positive, conflicts([]string{"hardwareprofile"})),
	"backups":                attribute(optional, forceNew, list(link("backuppolicy_vdc"))),
	"bootstrap":              attribute(optional, forceNew, text),
	"deploy":                 attribute(optional, forceNew, boolean),
	"disks":                  attribute(optional, forceNew, list(href)),
	"fws":                    attribute(optional, forceNew, list(link("firewall"))),
	"hardwareprofile":        attribute(optional, forceNew, link("hardwareprofile"), conflicts([]string{"cpu", "ram"})),
	"label":                  attribute(optional, forceNew, text),
	"lbs":                    attribute(optional, forceNew, list(link("loadbalancer"))),
	"ips":                    attribute(optional, forceNew, list(link("virtualmachine_ip"))),
	"monitored":              attribute(optional, forceNew, boolean),
	"name":                   attribute(computed, forceNew, text),
	"ram":                    attribute(optional, forceNew, positive, conflicts([]string{"hardwareprofile"})),
	"variables":              attribute(optional, forceNew, hash(text)),
	"virtualappliance":       attribute(required, forceNew, link("virtualappliance")),
	"virtualmachinetemplate": attribute(required, forceNew, link("template")),
}

func vmNew(d *resourceData) core.Resource {
	variables := make(map[string]string)
	for key, value := range d.dict("variables") {
		variables[key] = value.(string)
	}
	return &abiquo.VirtualMachine{
		CPU:       d.integer("cpu"),
		RAM:       d.integer("ram"),
		Label:     d.string("label"),
		Monitored: d.boolean("monitored"),
		Variables: variables,
		DTO: core.NewDTO(
			d.link("hardwareprofile"),
			d.link("virtualmachinetemplate"),
		),
	}
}

func vmEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("virtualappliance")+"/virtualmachines", "virtualmachine")
}

func vmReconfigure(vm *abiquo.VirtualMachine, d *resourceData) (err error) {
	// Update metadata
	if bootstrap, ok := d.GetOk("bootstrap"); ok {
		if err = vm.SetMetadata(&abiquo.VirtualMachineMetadata{
			Metadata: abiquo.VirtualMachineMetadataFields{
				StartupScript: bootstrap.(string),
			},
		}); err != nil {
			return
		}
	}

	fwsList := d.slice("fws")
	lbsList := d.slice("lbs")
	ipsList := d.slice("ips")
	hdsList := d.slice("disks")
	bckList := d.slice("backups")
	reconfigure := len(hdsList)+len(fwsList)+len(lbsList)+len(ipsList) > 0
	if reconfigure {
		// CONFIGURE disks
		for _, d := range hdsList {
			disk := new(abiquo.HardDisk)
			if err = core.Read(hdLink(d.(string)), disk); err != nil {
				return
			}
			if err = vm.AttachDisk(disk); err != nil {
				return
			}
		}

		// CONFIGURE nics
		for _, ip := range ipsList {
			if err = vm.AttachNIC(ipLink(ip.(string))); err != nil {
				return
			}
		}

		// CONFIGURE fws
		for _, fw := range fwsList {
			fwLink := core.NewLinkType(fw.(string), "firewallpolicy")
			vm.Add(fwLink.SetRel("firewall"))
		}

		// CONFIGURE lbs
		for _, lb := range lbsList {
			lbLink := core.NewLinkType(lb.(string), "loadbalancer")
			vm.Add(lbLink.SetRel("loadbalancer"))
		}

		// CONFIGURE backup policies
		for _, bck := range bckList {
			vm.Backups = append(vm.Backups, abiquo.BackupPolicy{
				DTO: core.NewDTO(
					core.NewLinkType(bck.(string), "backuppolicy").SetRel("policy"),
				),
			})
		}

		err = vm.Reconfigure()
	}
	return
}

func vmCreate(d *resourceData, resource core.Resource) (err error) {
	vm := resource.(*abiquo.VirtualMachine)
	if err = vmReconfigure(vm, d); err != nil {
		vm.Delete()
		return
	}

	d.SetId(vm.URL())
	if d.boolean("deploy") {
		err = vm.Deploy()
	}

	return
}

func vmRead(d *resourceData, resource core.Resource) (err error) {
	vm := resource.(*abiquo.VirtualMachine)
	d.Set("label", vm.Label)
	d.Set("name", vm.Name)
	d.Set("variables", vm.Variables)
	d.Set("virtualappliance", vm.Rel("virtualappliance").URL())
	d.Set("virtualmachinetemplate", vm.Rel("virtualmachinetemplate").URL())
	if _, ok := d.GetOk("profile"); ok {
		d.Set("profile", vm.Rel("hardwareprofile").URL())
	} else {
		if _, ok := d.GetOk("cpu"); ok {
			d.Set("cpu", vm.CPU)
		}
		if _, ok := d.GetOk("ram"); ok {
			d.Set("ram", vm.RAM)
		}
	}
	return
}

func vmUpdate(rd *schema.ResourceData, m interface{}) (err error) {
	d := newResourceData(rd, "virtualmachine")
	vm := vmNew(d).(*abiquo.VirtualMachine)
	return vm.Reconfigure()
}

func vmDelete(rd *schema.ResourceData, m interface{}) (err error) {
	d := newResourceData(rd, "virtualmachine")
	resource := d.Walk()
	if resource == nil {
		return
	}

	// To prevent the VM undeploy/delete sequence from breaking the vapp/vm
	// dependency, we have to undeploy the VM first if deployed, and delete it
	// once the VM is not allocated
	vm := resource.(*abiquo.VirtualMachine)
	if vm.State == "ON" || vm.State == "OFF" {
		if err = vm.Undeploy(); err != nil {
			return
		}
		vm = d.Walk().(*abiquo.VirtualMachine)
	}

	if vm.State != "NOT_ALLOCATED" {
		return fmt.Errorf("the VM is %v. it will not be deleted", vm.State)
	}

	return core.Remove(vm)
}

var resourceVm = &schema.Resource{
	Schema: vmSchema,
	Read:   resourceRead(vmNew, vmRead, "virtualmachine"),
	Exists: resourceExists("virtualmachine"),
	Delete: vmDelete,
	Create: resourceCreate(vmNew, vmCreate, vmRead, vmEndpoint),
}
