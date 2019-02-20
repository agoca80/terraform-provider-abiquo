resource "abiquo_private" "test" {
  virtualdatacenter = "${data.abiquo_virtualdatacenter.test.id}"
  mask              = 24
  address           = "172.16.10.0"
  gateway           = "172.16.10.1"
  name              = "test private"
  dns1              = "8.8.8.8"
  dns2              = "4.4.4.4"
  suffix            = "test.bcn.com"
  # XXX workaround ABICLOUDPREMIUM-9660
  lifecycle = { ignore_changes = [ "dns1", "dns2" ] }
}

data "abiquo_virtualdatacenter" "test" { name = "${var.virtualdatacenter}" }
