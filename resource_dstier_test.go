package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var dsTierTestHelper = &testHelper{
	kind:  "abiquo_dstier",
	media: "datastoretier",
	config: `
	data "abiquo_datacenter" "test" { name = "datacenter 1" }

	resource "abiquo_dstier" "test" {
	  datacenter  = "${data.abiquo_datacenter.test.id}"
		description = "required description"
		enabled     = true
	  name        = "testAccDSTierBasic"
	  policy      = "PERFORMANCE"
	}
	`,
}

func TestAccDSTier_update(t *testing.T) {
	resource.Test(t, dsTierTestHelper.updateCase(t))
}
