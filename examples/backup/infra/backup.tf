# Provider configuration: Get these from the environment
variable "endpoint" { }
variable "username" { }
variable "password" { }

provider "abiquo" {
  endpoint       = "${var.endpoint}"
  username       = "${var.username}"
  password       = "${var.password}"
}

data "abiquo_datacenter" "datacenter" { name = "datacenter 1" }
data "abiquo_location"   "location"   { name = "datacenter 1" }
data "abiquo_template"   "template"   { name = "tests" }

resource "abiquo_backup" "policy0" {
  datacenter     = "${data.abiquo_datacenter.datacenter.id}"
  name           = "terraform policy0"
  code           = "terraform policy0"
  description    = "optional"
  configurations = [
    { type    = "COMPLETE",   subtype = "HOURLY",         time = "2" },
    { type    = "SNAPSHOT",   subtype = "DAILY",          time = "01:02:03 +0400" },
    { type    = "FILESYSTEM", subtype = "MONTHLY",        time = "05:06:07 +0800" },
    { type    = "COMPLETE",   subtype = "WEEKLY_PLANNED", days = [ "wednesday", "monday", "tuesday", "thursday", "friday", "saturday", "sunday" ] }
  ]
}

resource "abiquo_backup" "policy1" {
  datacenter     = "${data.abiquo_datacenter.datacenter.id}"
  name           = "terraform policy1"
  code           = "terraform policy1"
  description    = "optional"
  configurations = [
    { type    = "COMPLETE",   subtype = "HOURLY",         time = "2" },
    { type    = "SNAPSHOT",   subtype = "DAILY",          time = "01:02:03 +0400" }
  ]
}

resource "abiquo_enterprise" "enterprise" {
  name = "backup"
}

resource "abiquo_limit" "limit" {
  enterprise = "${abiquo_enterprise.enterprise.id}"
  location   = "${data.abiquo_datacenter.datacenter.id}"
  backups    = [
    "${abiquo_backup.policy0.id}",
    "${abiquo_backup.policy1.id}"
  ]
}
