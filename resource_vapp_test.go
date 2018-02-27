package main

import (
	"fmt"
	"testing"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var testAccAbiquoVappBasic = `
data "abiquo_location"   "location"   { name = "datacenter 1" }
data "abiquo_enterprise" "enterprise" { name = "Abiquo" }

resource "abiquo_vdc" "vdc" {
	enterprise = "${data.abiquo_enterprise.enterprise.id}"
	location   = "${data.abiquo_location.location.id}"
  name       = "testAccAbiquoVappBasic"
	type       = "VMX_04"
}

resource "abiquo_vapp" "vapp" {
	virtualdatacenter = "${abiquo_vdc.vdc.id}"
  name              = "testAccAbiquoVappBasic"
}
`

func TestAccVapp_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckVappDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAbiquoVappBasic,
				Check: resource.ComposeTestCheckFunc(
					testCheckVappExists("abiquo_vapp.vapp"),
				),
			},
		},
	})
}

func testCheckVappDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abiquo_vapp" {
			continue
		}
		vapp := new(abiquo.VirtualAppliance)
		href := rs.Primary.Attributes["id"]
		endpoint := core.NewLinkType(href, "virtualappliance")
		if err := core.Read(endpoint, vapp); err == nil {
			return fmt.Errorf("virtual appliance %q still exists", vapp.Name)
		}
	}
	return nil
}

func testCheckVappExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("virtual appliance %q not found", name)
		}

		href := rs.Primary.Attributes["id"]
		endpoint := core.NewLinkType(href, "virtualappliance")
		return core.Read(endpoint, nil)
	}
}
