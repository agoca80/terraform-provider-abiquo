package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var sgTestHelper = &testHelper{
	kind:  "abiquo_sg",
	media: "scalinggroup",
	config: `
	data "abiquo_location"   "test" { name = "datacenter 1" }
	data "abiquo_enterprise" "test" { name = "Abiquo" }

	resource "abiquo_vdc" "test" {
	  enterprise = "${data.abiquo_enterprise.test.id}"
	  location   = "${data.abiquo_location.test.id}"
	  name       = "tf example sg"
	  type       = "KVM"
	}

	resource "abiquo_vapp" "test" {
	  virtualdatacenter = "${abiquo_vdc.test.id}"  
	  name = "test"
	}

	data "abiquo_template" "template" { name = "tests" }

	# Scaling group master instance
	resource "abiquo_vm" "test" {
	  virtualappliance       = "${abiquo_vapp.test.id}"
	  virtualmachinetemplate = "${data.abiquo_template.template.id}"
	  label                  = "test"
	}

	resource "abiquo_sg" "test" {
	  mastervirtualmachine = "${abiquo_vm.test.id}"
	  virtualappliance     = "${abiquo_vapp.test.id}"

	  name      = "test"
	  cooldown  = 60
	  min       = 2
	  max       = 4
	  scale_in  = [ { numberofinstances = 1 } ]
	  scale_out = [ { numberofinstances = 1 } ]
	}
	`,
}

func TestAccSG_update(t *testing.T) {
	resource.Test(t, sgTestHelper.updateCase(t))
}
