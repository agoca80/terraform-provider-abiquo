package main

import (
	"fmt"
	"testing"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var testAccBackupBasic = `
data "abiquo_datacenter"   "datacenter"   { name = "datacenter 1" }

resource "abiquo_backup" "backup" {
	# endpoint
	datacenter     = "${data.abiquo_datacenter.datacenter.id}"
	
	code           = "testAccBackupBasic (required)"
  name           = "testAccBackupBasic (required)"
	description    = "testAccBackupBasic (optional)"
	configurations = [ 
		{ type = "COMPLETE", subtype = "HOURLY", time = "2" }
	]
}

`

func TestAccBackup_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckBackupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccBackupBasic,
				Check: resource.ComposeTestCheckFunc(
					testCheckBackupExists("abiquo_backup.backup"),
				),
			},
		},
	})
}

func testCheckBackupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abiquo_backup" {
			continue
		}
		bck := new(abiquo.BackupPolicy)
		href := rs.Primary.Attributes["id"]
		endpoint := core.NewLinkType(href, "backuppolicy")
		if err := core.Read(endpoint, bck); err == nil {
			return fmt.Errorf("backup policy %q still exists", bck.Code)
		}
	}
	return nil
}

func testCheckBackupExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("backup policy %q not found", name)
		}

		href := rs.Primary.Attributes["id"]
		endpoint := core.NewLinkType(href, "backuppolicy")
		return core.Read(endpoint, nil)
	}
}
