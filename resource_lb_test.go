package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var lbTestHelper = &testHelper{
	kind:  "abiquo_lb",
	media: "loadbalancer",
	config: `
	data "abiquo_location"   "test" { name = "datacenter 1" }
	data "abiquo_enterprise" "test" { name = "Abiquo" }

	resource "abiquo_vdc" "test" {
	  enterprise = "${data.abiquo_enterprise.test.id}"
	  location   = "${data.abiquo_location.test.id}"
	  name       = "testAccLB"
	  type       = "KVM"
	}

	resource "abiquo_private" "test" {
	  virtualdatacenter = "${abiquo_vdc.test.id}"

	  # XXX workaround ABICLOUDPREMIUM-9660
	  lifecycle = { ignore_changes = [ "dns1", "dns2" ] }

	  mask    = 24
	  address = "172.16.27.0"
	  gateway = "172.16.27.1"
	  name    = "testAccLB"
	  dns1    = "8.8.8.8"
	  dns2    = "4.4.4.4"
	  suffix  = "test.abiquo.com"
	}

	resource "abiquo_lb" "test" {
	  virtualdatacenter = "${abiquo_vdc.test.id}"
	  privatenetwork    = "${abiquo_private.test.id}"

	  name         = "testAccLB"
	  algorithm    = "ROUND_ROBIN"
	  routingrules = [
	    { protocolin = "HTTP" , protocolout = "HTTP" , portin = 80 , portout = 80 }
	  ]
	} 
	`,
}

func TestAccLB_update(t *testing.T) {
	resource.Test(t, lbTestHelper.updateCase(t))
}
