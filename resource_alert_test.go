package main

import (
	"fmt"
	"testing"

	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var testAccAlertBasic = `
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
	name              = "testAccAlertBasic"
}

resource "abiquo_vm" "vm" {
	cpu                    = 1
	deploy                 = false
	ram                    = 64
	label                  = "testAccAlertBasic"
	virtualappliance       = "${abiquo_vapp.vapp.id}"
	virtualmachinetemplate = "${data.abiquo_template.template.id}"
}

resource "abiquo_alarm" "alarm" {
  target      = "${abiquo_vm.vm.id}"
  name        = "testAccAlertBasic"
  metric      = "cpu_time"
  period      = 60
  evaluations = 3
  statistic   = "average"
  formula     = "lessthan"
  threshold   = 10000
}

resource "abiquo_alert" "alert" {
  virtualappliance = "${abiquo_vapp.vapp.id}"
  name        = "testAccAlertBasic"
  description = "testAccAlertBasic"
  
  alarms = [
    "${abiquo_alarm.alarm.id}"
  ]
}
`

func TestAccAlert_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAlertDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAlertBasic,
				Check: resource.ComposeTestCheckFunc(
					testCheckAlertExists("abiquo_alert.alert"),
				),
			},
		},
	})
}

func testCheckAlertDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abiquo_alert" {
			continue
		}
		href := rs.Primary.Attributes["id"]
		endpoint := core.NewLinkType(href, "alert")
		if err := core.Read(endpoint, nil); err == nil {
			return fmt.Errorf("alert %q still exists", rs.Primary.Attributes["name"])
		}
	}
	return nil
}

func testCheckAlertExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("alert %q not found", name)
		}
		href := rs.Primary.Attributes["id"]
		endpoint := core.NewLinkType(href, "alert")
		return core.Read(endpoint, nil)
	}
}
