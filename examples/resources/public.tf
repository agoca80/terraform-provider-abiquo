resource "abiquo_public" "example" {
  datacenter         = "${data.abiquo_datacenter.datacenter.id}"
  networkservicetype = "${data.abiquo_datacenter.datacenter.id}/networkservicetypes/4"

  # XXX workaround ABICLOUDPREMIUM-9660
  lifecycle = { ignore_changes = [ "dns1", "dns2" ] }

  tag     = 3331
  mask    = 24
  address = "172.16.178.0"
  gateway = "172.16.178.1"
  name    = "terraform example"
  dns1    = "4.4.4.4"
  dns2    = "8.8.8.8"
  suffix  = "public.com"
}
