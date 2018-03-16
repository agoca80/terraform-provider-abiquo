package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var scopeTestHelper = &testHelper{
	kind:  "abiquo_scope",
	media: "scope",
	config: `
	data "abiquo_enterprise" "test" { name = "Abiquo" }
	data "abiquo_datacenter" "test" { name = "datacenter 1" }

	resource "abiquo_scope" "test" {
	  name        = "testAccAbiquoVappBasic"
	  datacenters = [ "${data.abiquo_datacenter.test.id}" ]
	  enterprises = [ "${data.abiquo_enterprise.test.id}" ]
	}
	`,
}

func TestAccScope_update(t *testing.T) {
	resource.Test(t, scopeTestHelper.updateCase(t))
}
