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
