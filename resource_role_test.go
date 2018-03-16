package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var roleTestHelper = &testHelper{
	kind:  "abiquo_role",
	media: "role",
	config: `
	resource "abiquo_role" "test" {
	  name = "test"
	  privileges = [
	    "APPLIB_UPLOAD_IMAGE",
	    "VAPP_CREATE_STATEFUL",
	    "VDC_MANAGE_VAPP",
	    "VM_ACTION_PLAN_MANAGE",
	  ]
	}
	`,
}

func TestAccRole_update(t *testing.T) {
	resource.Test(t, roleTestHelper.updateCase(t))
}
