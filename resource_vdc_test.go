package main

import (
	"fmt"
	"testing"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var testAccAbiquoVDCBasic = `
data "abiquo_location"   "location"   { name = "datacenter 1" }
data "abiquo_enterprise" "enterprise" { name = "Abiquo" }

resource "abiquo_vdc" "vdc" {
	enterprise = "${data.abiquo_enterprise.enterprise.id}"
	location   = "${data.abiquo_location.location.id}"
  name       = "testAccAbiquoVDCBasic"
	type       = "VMX_04"
}
`

func TestAccVdc_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckVDCDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAbiquoVDCBasic,
				Check: resource.ComposeTestCheckFunc(
					testCheckVDCExists("abiquo_vdc.vdc"),
				),
			},
		},
	})
}

func testCheckVDCDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abiquo_vdc" {
			continue
		}
		vdc := new(abiquo.VirtualDatacenter)
		href := rs.Primary.Attributes["id"]
		endpoint := core.NewLinkType(href, "virtualdatacenter")
		if err := core.Read(endpoint, vdc); err == nil {
			return fmt.Errorf("virtual datacenter %q still exists", vdc.Name)
		}
	}
	return nil
}

func testCheckVDCExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("virtual datacenter %q not found", name)
		}

		href := rs.Primary.Attributes["id"]
		endpoint := core.NewLinkType(href, "virtualdatacenter")
		return core.Read(endpoint, nil)
	}
}
