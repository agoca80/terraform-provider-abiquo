resource "abiquo_vdc" "vdc" {
  enterprise = "${data.abiquo_enterprise.enterprise.id}"
  location   = "${data.abiquo_location.location.id}"

  name       = "terraform vdc"
  type       = "VMX_04"
}
