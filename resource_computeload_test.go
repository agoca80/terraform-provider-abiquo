package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var computeTestHelper = &testHelper{
	kind:  "abiquo_computeload",
	media: "machineloadrule",
	config: `
	data "abiquo_datacenter" "test" { name = "datacenter 1" }

	resource "abiquo_computeload" "test" {
	  cpuload    = "1000"
	  ramload    = "95"
	  target     = "${data.abiquo_datacenter.test.id}"
	}
	`,
}

func TestAccCompute_update(t *testing.T) {
	resource.Test(t, computeTestHelper.updateCase(t))
}
