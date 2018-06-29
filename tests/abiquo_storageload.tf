data "abiquo_datacenter" "test" { name = "datacenter 1" }

resource "abiquo_storageload" "test" {
  load   = "95"
  target = "${data.abiquo_datacenter.test.id}"
}
