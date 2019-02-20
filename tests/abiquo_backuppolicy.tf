resource "abiquo_backuppolicy" "test" {
  datacenter     = "${data.abiquo_datacenter.test.id}"
  code           = "test backuppolicy"
# description    = "test backuppolicy"
  name           = "test backuppolicy"
  configurations = [
    { type = "COMPLETE", subtype = "HOURLY", time = "2" }
  ]
}

data "abiquo_datacenter" "test" { name = "${var.datacenter}" }
