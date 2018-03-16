package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var limitTestHelper = &testHelper{
	kind:  "abiquo_limit",
	media: "limit",
	config: `
	data "abiquo_datacenter" "test" { name = "datacenter 1" }
	data "abiquo_dstier"     "test" { 
		name       = "Default Tier"
		datacenter = "${data.abiquo_datacenter.test.id}"
	}

	resource "abiquo_enterprise" "test" {
	   name = "testAccLimit"
	   cpusoft  = 6    , cpuhard  = 8
	   ramsoft  = 8192 , ramhard  = 16384
	   vlansoft = 2    , vlanhard = 4
	}

	resource "abiquo_limit" "test" {
	  enterprise = "${abiquo_enterprise.test.id}"
	  location   = "${data.abiquo_datacenter.test.id}"
		dstiers    = [
			"${data.abiquo_dstier.test.id}"
		]
	}
	`,
}

func TestAccLimit_update(t *testing.T) {
	resource.Test(t, limitTestHelper.updateCase(t))
}
