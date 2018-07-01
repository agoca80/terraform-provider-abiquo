data "abiquo_datacenter" "test" {
  name = "datacenter 1"
}

data "abiquo_dstier" "test" {
  name       = "Default Tier"
  datacenter = "${data.abiquo_datacenter.test.id}"
}

resource "abiquo_hp" "test" {
  active     = true
  name       = "test limit"
  cpu        = 16
  ram        = 64
  datacenter = "${data.abiquo_datacenter.test.id}"
}

resource "abiquo_backup" "test" {
  datacenter  = "${data.abiquo_datacenter.test.id}"
  code        = "test limit"
  name        = "test limit"
  description = "test limit"

  configurations = [
    {
      type    = "COMPLETE"
      subtype = "HOURLY"
      time    = "2"
    },
  ]
}

resource "abiquo_enterprise" "test" {
  name = "test limit"
}

resource "abiquo_limit" "test" {
  enterprise = "${abiquo_enterprise.test.id}"
  location   = "${data.abiquo_datacenter.test.id}"
  dstiers    = ["${data.abiquo_dstier.test.id}"]
  backups    = ["${abiquo_backup.test.id}"]
  hwprofiles = ["${abiquo_hp.test.id}"]
}
