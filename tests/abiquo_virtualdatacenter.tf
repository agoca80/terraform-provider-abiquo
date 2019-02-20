resource "abiquo_virtualdatacenter" "test" {
  enterprise = "${data.abiquo_enterprise.enterprise.id}"
  location   = "${data.abiquo_location.location.id}"
  name       = "test vdc"
  net_mask   = "24"
  net_name   = "test vdc"
  net_address = "192.168.0.0"
  net_gateway = "192.168.0.1"
  type       = "KVM"
}

data "abiquo_location"   "location"   { name = "${var.datacenter}" }
data "abiquo_enterprise" "enterprise" { name = "Abiquo" }
