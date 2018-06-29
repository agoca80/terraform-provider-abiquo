data "abiquo_vdc"        "test"       { name = "tests" }

resource "abiquo_vol" "test" {
  tier               = "${abiquo_vdc.test.id}/tiers/1"
  virtualdatacenter  = "${abiquo_vdc.test.id}"

  type = "SCSI"
  name = "test"
  size = 32
}
