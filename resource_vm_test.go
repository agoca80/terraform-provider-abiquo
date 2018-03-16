package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var vmTestHelper = &testHelper{
	kind:  "abiquo_vm",
	media: "virtualmachine",
	config: `
	data "abiquo_enterprise" "test" { name = "Abiquo" }
	data "abiquo_location"   "test"   { name = "datacenter 1" }
	data "abiquo_template"   "test"   { name = "tests" }

	resource "abiquo_vdc" "test" {
		enterprise = "${data.abiquo_enterprise.test.id}"
		location   = "${data.abiquo_location.test.id}"
		name       = "testAccVMBasic"
		type       = "KVM"
	}

	resource "abiquo_fw" "test" {
		virtualdatacenter = "${abiquo_vdc.test.id}"

		description = "testAccVMBasic"
		name        = "testAccVMBasic"

		# XXX workaround ABICLOUDPREMIUM-9668
		rules = [
			{ protocol = "TCP", fromport = 22, toport = 22, sources = ["0.0.0.0/0"] }
		]
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

		name         = "testAccVMBasic"
		algorithm    = "ROUND_ROBIN"
		routingrules = [
			{ protocolin = "HTTP" , protocolout = "HTTP" , portin = 80 , portout = 80 }
		]
	}

	resource "abiquo_vapp" "test" {
		virtualdatacenter = "${abiquo_vdc.test.id}"
		name              = "testAccVMBasic"
	}

	resource "abiquo_vm" "test" {
		deploy                 = false
		backups                = [ ]
		cpu                    = 1
		ram                    = 64
		label                  = "testAccVMBasic"
		virtualappliance       = "${abiquo_vapp.test.id}"
		virtualmachinetemplate = "${data.abiquo_template.test.id}"

		lbs = [ "${abiquo_lb.test.id}" ]
		fws = [ "${abiquo_fw.test.id}" ]

		variables = {
			name1 = "value1"
			name2 = "value2"
		}

		bootstrap = <<EOF
	#!/bin/sh
	exit 0
	EOF
	}
	`,
}

func TestAccVM_update(t *testing.T) {
	resource.Test(t, vmTestHelper.updateCase(t))
}
