package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var userTestHelper = &testHelper{
	kind:  "abiquo_user",
	media: "user",
	config: `
  data "abiquo_enterprise" "test" { name = "Abiquo" }
	data "abiquo_role"       "test" { name = "CLOUD_ADMIN" }
	
	resource "abiquo_user" "test" {
	  enterprise = "${data.abiquo_enterprise.test.id}"
	  role       = "${data.abiquo_role.test.id}"
	  active     = true
	  name       = "test"
	  surname    = "test"
	  nick       = "test"
	  email      = "test@test.com"
	}
  `,
}

func TestAccUser_update(t *testing.T) {
	resource.Test(t, userTestHelper.updateCase(t))
}
