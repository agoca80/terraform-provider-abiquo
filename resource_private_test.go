package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var privateTestHelper = &testHelper{
	kind:  "abiquo_private",
	media: "vlan",
	config: `
	data "abiquo_location"   "test" { name = "datacenter 1" }
	data "abiquo_enterprise" "test" { name = "Abiquo" }

	resource "abiquo_vdc" "test" {
	  enterprise = "${data.abiquo_enterprise.test.id}"
	  location   = "${data.abiquo_location.test.id}"

	  name       = "testAccPrivate"
	  type       = "KVM"
	}

	resource "abiquo_private" "test" {
	  virtualdatacenter = "${abiquo_vdc.test.id}"

	  # XXX workaround ABICLOUDPREMIUM-9660
	  lifecycle = { ignore_changes = [ "dns1", "dns2" ] }

	  mask    = 24
	  address = "172.16.10.0"
	  gateway = "172.16.10.1"
	  name    = "testAccPrivate"
	  dns1    = "8.8.8.8"
	  dns2    = "4.4.4.4"
	  suffix  = "test.bcn.com"
	}
	`,
}

func TestAccPrivate_update(t *testing.T) {
	resource.Test(t, privateTestHelper.updateCase(t))
}
