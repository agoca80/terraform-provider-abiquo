package main

import (
	"fmt"
	"testing"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var testAccAbiquoRoleBasic = `
resource "abiquo_role" "role" {
  name = "testAccAbiquoRoleBasic"
  privileges = [
    "APPLIB_UPLOAD_IMAGE",
    "VAPP_CREATE_STATEFUL",
    "VDC_MANAGE_VAPP",
    "VM_ACTION_PLAN_MANAGE",
  ]
}
`

func TestAccRole_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckEnterpriseDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAbiquoRoleBasic,
				Check: resource.ComposeTestCheckFunc(
					testCheckRoleExists("abiquo_role.role"),
				),
			},
		},
	})
}

func testCheckRoleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abiquo_role" {
			continue
		}
		role := new(abiquo.Role)
		href := rs.Primary.Attributes["id"]
		endpoint := core.NewLinkType(href, "role")
		if err := core.Read(endpoint, role); err == nil {
			return fmt.Errorf("role %q still exists", role.Name)
		}
	}
	return nil
}

func testCheckRoleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("role %q not found", name)
		}

		href := rs.Primary.Attributes["id"]
		endpoint := core.NewLinkType(href, "role")
		return core.Read(endpoint, nil)
	}
}
