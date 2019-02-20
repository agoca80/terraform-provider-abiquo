resource "abiquo_firewallpolicy" "test" {
  device            = "${data.abiquo_virtualdatacenter.test.device}"
  virtualdatacenter = "${data.abiquo_virtualdatacenter.test.id}"
  description       = "test fw"
  name              = "test fw"
  # XXX workaround ABICLOUDPREMIUM-9668
  rules = [
    { protocol = "TCP", fromport = 22, toport = 22, sources = ["0.0.0.0/0"] },
    { protocol = "TCP", fromport = 80, toport = 80, sources = ["0.0.0.0/0"] },
    { protocol = "TCP", fromport = 44, toport = 44, sources = ["0.0.0.0/0"] }
  ]
}

data "abiquo_location"   "test" { name = "${var.datacenter}" }
data "abiquo_enterprise" "test" { name = "Abiquo" }
data "abiquo_virtualdatacenter"        "test" { name = "${var.virtualdatacenter}" }
