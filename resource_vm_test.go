package main

import (
	"fmt"
	"testing"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var testAccAbiquoVMBasic = `
data "abiquo_enterprise" "enterprise" { name = "Abiquo" }
data "abiquo_location"   "location"   { name = "datacenter 1" }
data "abiquo_template"   "template"   { name = "tests" }

resource "abiquo_vdc" "vdc" {
	enterprise = "${data.abiquo_enterprise.enterprise.id}"
	location   = "${data.abiquo_location.location.id}"
	name       = "testAccVMBasic"
	type       = "VMX_04"
}

resource "abiquo_vapp" "vapp" {
	virtualdatacenter = "${abiquo_vdc.vdc.id}"
	name              = "testAccVMBasic"
}

resource "abiquo_vm" "vm" {
	backups                = [ ]
	cpu                    = 1
	ram                    = 64
	label                  = "testAccVMBasic"
	virtualappliance       = "${abiquo_vapp.vapp.id}"
	virtualmachinetemplate = "${data.abiquo_template.template.id}"
}
`

func TestAccVM_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckVMDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAbiquoVMBasic,
				Check: resource.ComposeTestCheckFunc(
					testCheckVMExists("abiquo_vm.vm"),
				),
			},
		},
	})
}

func testCheckVMDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abiquo_vm" {
			continue
		}
		vm := new(abiquo.VirtualMachine)
		href := rs.Primary.Attributes["id"]
		endpoint := core.NewLinkType(href, "virtualmachine")
		if err := core.Read(endpoint, vm); err == nil {
			return fmt.Errorf("VM %q still exists", vm.Label)
		}
	}
	return nil
}

func testCheckVMExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("VM %q not found", name)
		}

		href := rs.Primary.Attributes["id"]
		endpoint := core.NewLinkType(href, "virtualmachine")
		return core.Read(endpoint, nil)
	}
}
