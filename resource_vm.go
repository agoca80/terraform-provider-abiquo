package main

import (
	"github.com/abiquo/opal/abiquo"
	"github.com/abiquo/opal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var vmResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"cpu":                    Conflicts([]string{"hardwareprofile"}).Renew().Number(),
		"ram":                    Conflicts([]string{"hardwareprofile"}).Renew().Number(),
		"hardwareprofile":        Conflicts([]string{"cpu", "ram"}).Renew().Link(),
		"bootstrap":              Optional().Renew().String(),
		"disks":                  Optional().Renew().Links(),
		"fws":                    Optional().Renew().Links(),
		"label":                  Optional().Renew().String(),
		"lbs":                    Optional().Renew().Links(),
		"monitored":              Optional().Renew().Bool(),
		"ips":                    Optional().Renew().Links(),
		"name":                   Computed().Renew().String(),
		"variables":              New(optional(&schema.Schema{Type: schema.TypeMap, Elem: fieldString()})),
		"virtualappliance":       Required().Renew().Link(),
		"virtualmachinetemplate": Required().Renew().Link(),
	},
	Read:   resourceRead(vmNew, vmRead, "virtualmachine"),
	Exists: resourceExists("virtualmachine"),
	Delete: resourceDelete,
	Create: vmCreate,
}

func vmNew(d *resourceData) core.Resource {
	variables := make(map[string]string)
	for key, value := range d.dict("variables") {
		variables[key] = value.(string)
	}
	return &abiquo.VirtualMachine{
		CPU:       d.int("cpu"),
		RAM:       d.int("ram"),
		Label:     d.string("label"),
		Monitored: d.bool("monitored"),
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

	// BEGIN reconfigure
	fwsList := d.slice("fws")
	lbsList := d.slice("lbs")
	ipsList := d.slice("ips")
	hdsList := d.slice("disks")
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

		err = vm.Reconfigure()
	}
	return
}

func vmCreate(rd *schema.ResourceData, m interface{}) (err error) {
	d := newResourceData(rd, "")
	vm := vmNew(d).(*abiquo.VirtualMachine)
	if err = core.Create(vmEndpoint(d), vm); err != nil {
		return
	}
	d.SetId(vm.URL())

	if err = vmReconfigure(vm, d); err == nil {
		err = vm.Deploy()
	}
	return
}

func vmRead(d *resourceData, resource core.Resource) (err error) {
	vm := resource.(*abiquo.VirtualMachine)
	d.Set("label", vm.Label)
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
