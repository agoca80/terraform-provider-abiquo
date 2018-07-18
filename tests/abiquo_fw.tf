resource "abiquo_fw" "test" {
  virtualdatacenter = "${data.abiquo_vdc.test.id}"

  description = "test fw"
  name        = "test fw"

  # XXX workaround ABICLOUDPREMIUM-9668
  rules = [
    { protocol = "TCP", fromport = 22, toport = 22, sources = ["0.0.0.0/0"] },
    { protocol = "TCP", fromport = 80, toport = 80, sources = ["0.0.0.0/0"] },
    { protocol = "TCP", fromport = 443, toport = 443, sources = ["0.0.0.0/0"] }
  ]
}

data "abiquo_location"   "test" { name = "datacenter 1" }
data "abiquo_enterprise" "test" { name = "Abiquo" }
data "abiquo_vdc"        "test" { name = "tests" }
