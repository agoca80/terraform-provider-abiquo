package main

import (
	"fmt"
	"testing"

	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var testAccPlanBasic = `
data "abiquo_enterprise" "enterprise" { name = "Abiquo" }
data "abiquo_location"   "location"   { name = "datacenter 1" }
data "abiquo_template"   "template"   { name = "tests" }

resource "abiquo_vdc" "vdc" {
	enterprise = "${data.abiquo_enterprise.enterprise.id}"
	location   = "${data.abiquo_location.location.id}"
	name       = "testAccPlanBasic"
	type       = "VMX_04"
}

resource "abiquo_vapp" "vapp" {
	virtualdatacenter = "${abiquo_vdc.vdc.id}"
	name              = "testAccPlanBasic"
}

resource "abiquo_vm" "vm" {
	cpu                    = 1
	deploy                 = false
	ram                    = 64
	label                  = "testAccPlanBasic"
	virtualappliance       = "${abiquo_vapp.vapp.id}"
	virtualmachinetemplate = "${data.abiquo_template.template.id}"
}

resource "abiquo_plan" "plan" {
	virtualmachine = "${abiquo_vm.vm.id}"
	description    = "testAccPlanBasic"
	name           = "testAccPlanBasic"
	entries        = [
		{	parameter = "", parametertype = "None",	type = "UNDEPLOY" }
	]
}
`

func TestAccPlan_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckPlanDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccPlanBasic,
				Check: resource.ComposeTestCheckFunc(
					testCheckPlanExists("abiquo_plan.plan"),
				),
			},
		},
	})
}

func testCheckPlanDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abiquo_plan" {
			continue
		}
		href := rs.Primary.Attributes["id"]
		endpoint := core.NewLinkType(href, "virtualmachineactionplan")
		if err := core.Read(endpoint, nil); err == nil {
			return fmt.Errorf("testCheckPlanDestroy: action plan %q still exists", rs.Primary.Attributes["name"])
		}
	}
	return nil
}

func testCheckPlanExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("action plan %q not found", name)
		}
		href := rs.Primary.Attributes["id"]
		endpoint := core.NewLinkType(href, "virtualmachineactionplan")
		return core.Read(endpoint, nil)
	}
}
