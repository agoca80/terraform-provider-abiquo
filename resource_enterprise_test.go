package main

import (
	"fmt"
	"testing"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var testAccEnterpriseBasic = `
resource "abiquo_enterprise" "enterprise" {
    name       = "testAccEnterpriseBasic"
		properties = {
			"property0" = "value0"
			"property1" = "value1"
		}
}
`

func TestAccEnterprise_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckEnterpriseDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccEnterpriseBasic,
				Check: resource.ComposeTestCheckFunc(
					testCheckEnterpriseExists("abiquo_enterprise.enterprise"),
				),
			},
		},
	})
}

func testCheckEnterpriseDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abiquo_enterprise" {
			continue
		}
		enterprise := new(abiquo.Enterprise)
		href := rs.Primary.Attributes["id"]
		endpoint := core.NewLinkType(href, "enterprise")
		if err := core.Read(endpoint, enterprise); err == nil {
			return fmt.Errorf("enterprise %q still exists", enterprise.Name)
		}
	}
	return nil
}

func testCheckEnterpriseExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("enterprise %q not found", name)
		}

		href := rs.Primary.Attributes["id"]
		endpoint := core.NewLinkType(href, "enterprise")
		return core.Read(endpoint, nil)
	}
}
