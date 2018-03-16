package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var storageLoadTestHelper = &testHelper{
	kind:  "abiquo_storageload",
	media: "datastoreloadrule",
	config: `
	data "abiquo_datacenter" "test" { name = "datacenter 1" }

	resource "abiquo_storageload" "test" {
	  load   = "95"
	  target = "${data.abiquo_datacenter.test.id}"
	}`,
}

func TestAccStorageLoad_update(t *testing.T) {
	resource.Test(t, storageLoadTestHelper.updateCase(t))
}
