resource "abiquo_rack" "test" {
  name        = "testAccRackBasic"
  vlanmin     = 1000
  vlanmax     = 1999
  description = "kvm"
  datacenter  = "${data.abiquo_datacenter.test.id}"
}

data "abiquo_datacenter" "test" { name = "datacenter 1" }
