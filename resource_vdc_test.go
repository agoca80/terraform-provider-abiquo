package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var vdcTestHelper = &testHelper{
	kind:  "abiquo_vdc",
	media: "virtualdatacenter",
	config: `
	data "abiquo_location"   "location"   { name = "datacenter 1" }
	data "abiquo_enterprise" "enterprise" { name = "Abiquo" }

	resource "abiquo_vdc" "test" {
		enterprise = "${data.abiquo_enterprise.enterprise.id}"
		location   = "${data.abiquo_location.location.id}"
	  name       = "testAccAbiquoVDCBasic"
		type       = "KVM"
	}
	`,
}

func TestAccVdc_update(t *testing.T) {
	resource.Test(t, vdcTestHelper.updateCase(t))
}
