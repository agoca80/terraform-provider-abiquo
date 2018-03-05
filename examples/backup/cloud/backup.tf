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
data "abiquo_enterprise" "enterprise" { name = "backup" }

data "abiquo_backup"     "policy" {
  count    = 2
  code     = "terraform policy${count.index}"
  location = "${data.abiquo_location.location.id}"
}

resource "abiquo_vdc" "backup" {
  enterprise = "${data.abiquo_enterprise.enterprise.id}"
  location   = "${data.abiquo_location.location.id}"

  name       = "terraform backup"
  type       = "VMX_04"
}

resource "abiquo_vapp" "backup" {
  virtualdatacenter = "${abiquo_vdc.backup.id}"
  name = "terraform backup"
}

resource "abiquo_vm" "backup" {
  deploy                 = false
  label                  = "terraform backup"
  virtualappliance       = "${abiquo_vapp.backup.id}"
  virtualmachinetemplate = "${data.abiquo_template.template.id}"
  
  backups        = [
    "${data.abiquo_backup.policy.0.id}",
    "${data.abiquo_backup.policy.1.id}"
  ]
}
