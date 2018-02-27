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
data "abiquo_vdc"      "vdc"      { name = "tests" }
data "abiquo_template" "template" { name = "tests" }

data "abiquo_vapp" "vapp" {
	virtualdatacenter = "${data.abiquo_vdc.vdc.id}"
	name              = "tests"
}

resource "abiquo_vm" "vm" {
	cpu                    = 1
	ram                    = 64
	label                  = "testAccAbiquoVMBasic"
	virtualappliance       = "${data.abiquo_vapp.vapp.id}"
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
