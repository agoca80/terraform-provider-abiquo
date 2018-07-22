resource "abiquo_volume" "test" {
  tier              = "${data.abiquo_tier.test.id}"
  virtualdatacenter = "${data.abiquo_virtualdatacenter.test.id}"

  type = "SCSI"
  name = "test vol"
  size = 32
}

data "abiquo_virtualdatacenter"  "test" { name = "tests" }
data "abiquo_tier" "test" {
  location = "${data.abiquo_virtualdatacenter.test.tiers}"
  name     = "Default Tier 1"
}
