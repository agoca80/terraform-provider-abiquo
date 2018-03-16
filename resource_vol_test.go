package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var volTestHelper = &testHelper{
	kind:  "abiquo_vol",
	media: "volume",
	config: `
  data "abiquo_enterprise" "enterprise" { name = "Abiquo" }
  data "abiquo_location"   "location"   { name = "datacenter 1" }

  resource "abiquo_vdc" "test" {
  	enterprise = "${data.abiquo_enterprise.enterprise.id}"
  	location   = "${data.abiquo_location.location.id}"
  	name       = "test"
  	type       = "KVM"
  }

  resource "abiquo_vol" "test" {
    tier               = "${abiquo_vdc.test.id}/tiers/1"
    virtualdatacenter  = "${abiquo_vdc.test.id}"

    type = "SCSI"
    name = "test"
    size = 32
  }
  `,
}

func TestAccVolume_update(t *testing.T) {
	resource.Test(t, volTestHelper.updateCase(t))
}
