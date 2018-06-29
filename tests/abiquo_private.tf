data "abiquo_vdc" "test" { name = "tests" }

resource "abiquo_private" "test" {
  virtualdatacenter = "${data.abiquo_vdc.test.id}"

  # XXX workaround ABICLOUDPREMIUM-9660
  lifecycle = { ignore_changes = [ "dns1", "dns2" ] }

  mask    = 24
  address = "172.16.10.0"
  gateway = "172.16.10.1"
  name    = "testAccPrivate"
  dns1    = "8.8.8.8"
  dns2    = "4.4.4.4"
  suffix  = "test.bcn.com"
}
