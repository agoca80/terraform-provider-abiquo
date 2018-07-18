resource "abiquo_public" "test" {
  datacenter         = "${data.abiquo_datacenter.test.id}"
  networkservicetype = "${data.abiquo_nst.test.id}"

  # XXX workaround ABICLOUDPREMIUM-9660
  lifecycle = { ignore_changes = [ "dns1", "dns2" ] }

  tag     = 3331
  mask    = 24
  address = "172.16.178.0"
  gateway = "172.16.178.1"
  name    = "test public"
  dns1    = "4.4.4.4"
  dns2    = "8.8.8.8"
  suffix  = "public.com"
}

data "abiquo_datacenter" "test" { name = "datacenter 1" }
data "abiquo_nst"        "test" {
  datacenter = "${data.abiquo_datacenter.test.id}"
  name       = "Service Network"
}
