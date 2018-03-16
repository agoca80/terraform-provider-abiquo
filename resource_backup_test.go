package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var backupTestHelper = &testHelper{
	kind:  "abiquo_backup",
	media: "backuppolicy",
	config: `
	data "abiquo_datacenter" "test" { name = "datacenter 1" }

	resource "abiquo_backup" "test" {
		# endpoint
		datacenter     = "${data.abiquo_datacenter.test.id}"

		code           = "testAccBackupBasic (required)"
	  name           = "testAccBackupBasic (required)"
		description    = "testAccBackupBasic (optional)"
		configurations = [ 
			{ type = "COMPLETE", subtype = "HOURLY", time = "2" }
		]
	}
	`,
}

func TestAccBackup_update(t *testing.T) {
	resource.Test(t, backupTestHelper.updateCase(t))
}
