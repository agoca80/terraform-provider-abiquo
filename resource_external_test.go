package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var externalTestHelper = &testHelper{
	kind:  "abiquo_external",
	media: "vlan",
	config: `
	data "abiquo_enterprise" "test" { name = "Abiquo" }
	data "abiquo_datacenter" "test" { name = "datacenter 1" }
	data "abiquo_nst"        "test" { 
		datacenter = "${data.abiquo_datacenter.test.id}"
		name       = "Service Network"
	}
	
	resource "abiquo_external" "test" {
	  enterprise         = "${data.abiquo_enterprise.test.id}"
	  datacenter         = "${data.abiquo_datacenter.test.id}"
	  networkservicetype = "${data.abiquo_nst.test.id}"

	  # XXX workaround ABICLOUDPREMIUM-9660
	  lifecycle = { ignore_changes = [ "dns1", "dns2" ] }

	  tag     = 1331
	  mask    = 24
	  address = "172.16.4.0"
	  gateway = "172.16.4.1"
	  name    = "testAccExternal"
	  dns1    = "4.4.4.4"
	  dns2    = "8.8.8.8"
	  suffix  = "external.test.abiquo.com"
	}
	`,
}

func TestAccExternal_update(t *testing.T) {
	resource.Test(t, externalTestHelper.updateCase(t))
}
