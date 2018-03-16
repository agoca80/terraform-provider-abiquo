package main

import (
	"fmt"
	"testing"

	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var testAccComputeLoadBasic = `
data "abiquo_datacenter" "datacenter" { name = "datacenter 1" }

resource "abiquo_computeload" "computeload" {
  cpuload    = "1000"
  ramload    = "95"
  target     = "${data.abiquo_datacenter.datacenter.id}"
}
`

func TestAccComputeLoad_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckComputeLoadDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeLoadBasic,
				Check: resource.ComposeTestCheckFunc(
					testCheckComputeLoadExists("abiquo_computeload.computeload"),
				),
			},
		},
	})
}

func testCheckComputeLoadDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abiquo_computeload" {
			continue
		}
		href := rs.Primary.Attributes["id"]
		endpoint := core.NewLinkType(href, "machineloadrule")
		if err := core.Read(endpoint, nil); err == nil {
			return fmt.Errorf("computeload %q still exists", rs.Primary.Attributes["name"])
		}
	}
	return nil
}

func testCheckComputeLoadExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("computeload %q not found", name)
		}
		href := rs.Primary.Attributes["id"]
		endpoint := core.NewLinkType(href, "machineloadrule")
		return core.Read(endpoint, nil)
	}
}
