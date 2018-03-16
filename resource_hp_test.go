package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var hpTestHelper = &testHelper{
	kind:  "abiquo_hp",
	media: "hardwareprofile",
	config: `
	data "abiquo_datacenter" "test" { name = "datacenter 1" }

	resource "abiquo_hp" "test" {
	  active = true
	  name = "testAccHP"
	  cpu  = 16
	  ram  = 64
	  datacenter = "${data.abiquo_datacenter.test.id}"
	}
	`,
}

func TestAccHP_update(t *testing.T) {
	resource.Test(t, hpTestHelper.updateCase(t))
}
