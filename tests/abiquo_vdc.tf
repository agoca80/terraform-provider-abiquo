data "abiquo_location"   "location"   { name = "datacenter 1" }
data "abiquo_enterprise" "enterprise" { name = "Abiquo" }

resource "abiquo_vdc" "test" {
  enterprise = "${data.abiquo_enterprise.enterprise.id}"
  location   = "${data.abiquo_location.location.id}"
  name       = "testAccAbiquoVDCBasic"
  type       = "KVM"
}
