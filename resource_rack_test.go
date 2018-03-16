package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var rackTestHelper = &testHelper{
	kind:  "abiquo_rack",
	media: "rack",
	config: `
	data "abiquo_datacenter" "test" { name = "datacenter 1" }

	resource "abiquo_rack" "test" {
	  name        = "testAccRackBasic"
	  vlanmin     = 1000
	  vlanmax     = 1999
	  description = "kvm"
	  datacenter  = "${data.abiquo_datacenter.test.id}"
	}
	`,
}

func TestAccRack_update(t *testing.T) {
	resource.Test(t, rackTestHelper.updateCase(t))
}
