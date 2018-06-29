data "abiquo_location"   "test" { name = "datacenter 1" }
data "abiquo_enterprise" "test" { name = "Abiquo" }
data "abiquo_vdc"        "test"   { name = "tests" }

resource "abiquo_private" "test" {
  virtualdatacenter = "${data.abiquo_vdc.test.id}"

  # XXX workaround ABICLOUDPREMIUM-9660
  lifecycle = { ignore_changes = [ "dns1", "dns2" ] }

  mask    = 24
  address = "172.16.27.0"
  gateway = "172.16.27.1"
  name    = "test lb"
  dns1    = "8.8.8.8"
  dns2    = "4.4.4.4"
  suffix  = "test.abiquo.com"
}

resource "abiquo_lb" "test" {
  virtualdatacenter = "${data.abiquo_vdc.test.id}"
  privatenetwork    = "${abiquo_private.test.id}"

  name         = "test lb"
  internal     = false
  algorithm    = "ROUND_ROBIN"
  routingrules = [
    { protocolin = "HTTP" , protocolout = "HTTP" , portin = 80 , portout = 80 }
  ]
}
