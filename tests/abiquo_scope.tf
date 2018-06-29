data "abiquo_enterprise" "test" { name = "Abiquo" }
data "abiquo_datacenter" "test" { name = "datacenter 1" }

resource "abiquo_scope" "test" {
  name        = "testAccAbiquoVappBasic"
  datacenters = [ "${data.abiquo_datacenter.test.id}" ]
  enterprises = [ "${data.abiquo_enterprise.test.id}" ]
}
