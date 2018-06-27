package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var vmSchema = map[string]*schema.Schema{
	"cpu": &schema.Schema{
		ConflictsWith: []string{"hardwareprofile"},
		ForceNew:      true,
		Optional:      true,
		Type:          schema.TypeInt,
	},
	"backups": &schema.Schema{
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validateURL,
		},
		ForceNew: true,
		Optional: true,
		Type:     schema.TypeList,
	},
	"bootstrap": &schema.Schema{
		ForceNew: true,
		Optional: true,
		Type:     schema.TypeString,
	},
	"deploy": &schema.Schema{
		ForceNew: true,
		Optional: true,
		Type:     schema.TypeBool,
	},
	"disks": &schema.Schema{
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validateURL,
		},
		ForceNew: true,
		Optional: true,
		Type:     schema.TypeList,
	},
	"fws": &schema.Schema{
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validateURL,
		},
		ForceNew: true,
		Optional: true,
		Type:     schema.TypeList,
	},
	"hardwareprofile": &schema.Schema{
		ConflictsWith: []string{"cpu", "ram"},
		ForceNew:      true,
		Optional:      true,
		Type:          schema.TypeString,
		ValidateFunc:  validateURL,
	},
	"label": &schema.Schema{
		ForceNew: true,
		Optional: true,
		Type:     schema.TypeString,
	},
	"lbs": &schema.Schema{
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validateURL,
		},
		ForceNew: true,
		Optional: true,
		Type:     schema.TypeList,
	},
	"monitored": &schema.Schema{
		ForceNew: true,
		Optional: true,
		Type:     schema.TypeBool,
	},
	"ips": &schema.Schema{
		Elem: &schema.Schema{
			Type: schema.TypeString,
			ValidateFunc: validateHref([]string{
				href["privateip"],
				href["externalip"],
				href["publicip"],
			}),
		},
		ForceNew: true,
		Optional: true,
		Type:     schema.TypeList,
	},
	"name": &schema.Schema{
		Computed: true,
		ForceNew: true,
		Type:     schema.TypeString,
	},
	"ram": &schema.Schema{
		ConflictsWith: []string{"hardwareprofile"},
		ForceNew:      true,
		Optional:      true,
		Type:          schema.TypeInt,
	},
	"variables": &schema.Schema{
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		ForceNew: true,
		Optional: true,
		Type:     schema.TypeMap,
	},
	"virtualappliance": &schema.Schema{
		ForceNew: true,
		Required: true,
		Type:     schema.TypeString,
		ValidateFunc: validateHref([]string{
			href["virtualappliance"],
		}),
	},
	"virtualmachinetemplate": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
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
	if d.bool("deploy") {
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
