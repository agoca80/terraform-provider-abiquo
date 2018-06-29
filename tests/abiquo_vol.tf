data "abiquo_vdc"        "test"       { name = "tests" }

resource "abiquo_vol" "test" {
  tier               = "${data.abiquo_vdc.test.id}/tiers/1"
  virtualdatacenter  = "${data.abiquo_vdc.test.id}"

  type = "SCSI"
  name = "test"
  size = 32
}
