package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var vappTestHelper = &testHelper{
	kind:  "abiquo_vapp",
	media: "virtualappliance",
	config: `
	data "abiquo_location"   "test"   { name = "datacenter 1" }
	data "abiquo_enterprise" "test" { name = "Abiquo" }

	resource "abiquo_vdc" "test" {
		enterprise = "${data.abiquo_enterprise.test.id}"
		location   = "${data.abiquo_location.test.id}"
	  name       = "testAccVappBasic"
		type       = "KVM"
	}

	resource "abiquo_vapp" "test" {
		virtualdatacenter = "${abiquo_vdc.test.id}"
	  name              = "testAccVappBasic"
	}
	`,
}

func TestAccVapp_update(t *testing.T) {
	resource.Test(t, vappTestHelper.updateCase(t))
}
