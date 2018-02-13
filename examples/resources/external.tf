resource "abiquo_external" "example" {
  enterprise         = "${abiquo_enterprise.example.id}"
  datacenter         = "${data.abiquo_datacenter.datacenter.id}"
  networkservicetype = "${data.abiquo_datacenter.datacenter.id}/networkservicetypes/4"

  # XXX workaround ABICLOUDPREMIUM-9660
  lifecycle = { ignore_changes = [ "dns1", "dns2" ] }

  tag     = 1331
  mask    = 24
  address = "172.16.4.0"
  gateway = "172.16.4.1"
  name    = "terraform example external"
  dns1    = "4.4.4.4"
  dns2    = "8.8.8.8"
  suffix  = "demo.com"
}
