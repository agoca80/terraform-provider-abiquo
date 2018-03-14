package main

import (
	"fmt"
	"testing"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var testAccDSTierBasic = `
data "abiquo_datacenter" "datacenter" { name = "datacenter 1" }

resource "abiquo_dstier" "dstier" {
  datacenter  = "${data.abiquo_datacenter.datacenter.id}"
	description = "required description"
	enabled     = true
  name        = "testAccDSTierBasic"
  policy      = "PERFORMANCE"
}

`

func TestAccDSTier_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckDSTierDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDSTierBasic,
				Check: resource.ComposeTestCheckFunc(
					testCheckDSTierExists("abiquo_dstier.dstier"),
				),
			},
		},
	})
}

func testCheckDSTierDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abiquo_dstier" {
			continue
		}
		dstier := new(abiquo.DatastoreTier)
		href := rs.Primary.Attributes["id"]
		endpoint := core.NewLinkType(href, "datastoretier")
		if err := core.Read(endpoint, dstier); err == nil {
			return fmt.Errorf("datastore tier %q still exists", dstier.Name)
		}
	}
	return nil
}

func testCheckDSTierExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("datastore tier %q not found", name)
		}

		href := rs.Primary.Attributes["id"]
		endpoint := core.NewLinkType(href, "datastoretier")
		return core.Read(endpoint, nil)
	}
}
