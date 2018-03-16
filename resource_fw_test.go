package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var fwTestHelper = &testHelper{
	kind:  "abiquo_fw",
	media: "firewallpolicy",
	config: `
	data "abiquo_location"   "test" { name = "datacenter 1" }
	data "abiquo_enterprise" "test" { name = "Abiquo" }

	resource "abiquo_vdc" "test" {
	  enterprise = "${data.abiquo_enterprise.test.id}"
	  location   = "${data.abiquo_location.test.id}"
	  name       = "testAccFW"
	  type       = "KVM"
	}

	resource "abiquo_fw" "test" {
	  virtualdatacenter = "${abiquo_vdc.test.id}"

	  description = "testAccFW"
	  name        = "testAccFW"

	  # XXX workaround ABICLOUDPREMIUM-9668
	  rules = [
	    { protocol = "TCP", fromport = 22, toport = 22, sources = ["0.0.0.0/0"] },
	    { protocol = "TCP", fromport = 80, toport = 80, sources = ["0.0.0.0/0"] },
	    { protocol = "TCP", fromport = 443, toport = 443, sources = ["0.0.0.0/0"] }
	  ]
	}
	`,
}

func TestAccFW_update(t *testing.T) {
	resource.Test(t, fwTestHelper.updateCase(t))
}
