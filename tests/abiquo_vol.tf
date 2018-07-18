resource "abiquo_vol" "test" {
  tier              = "${data.abiquo_tier.test.id}"
  virtualdatacenter = "${data.abiquo_vdc.test.id}"

  type = "SCSI"
  name = "test vol"
  size = 32
}

data "abiquo_vdc"  "test" { name = "tests" }
data "abiquo_tier" "test" {
  location = "${data.abiquo_vdc.test.tiers}"
  name     = "Default Tier 1"
}
