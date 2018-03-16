package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var enterpriseTestHelper = &testHelper{
	kind:  "abiquo_enterprise",
	media: "enterprise",
	config: `
	data "abiquo_datacenter" "test" { name = "datacenter 1" }

	resource "abiquo_enterprise" "test" {
	    name       = "testAccEnterpriseBasic"
			properties = {
				"property0" = "value0"
				"property1" = "value1"
			}
	}
	`,
}

func TestAccEnterprise_update(t *testing.T) {
	resource.Test(t, enterpriseTestHelper.updateCase(t))
}
