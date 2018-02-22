# Provider configuration: Get these from the environment
variable "endpoint" { }
variable "username" { }
variable "password" { }

provider "abiquo" {
  endpoint       = "${var.endpoint}"
  username       = "${var.username}"
  password       = "${var.password}"
}

variable "datacenter" { default = "datacenter 1" }

data "abiquo_datacenter" "datacenter" { name = "${var.datacenter}" }

data "abiquo_dstier"     "dstier"     { 
  name       = "Default Tier"
  datacenter = "${data.abiquo_datacenter.datacenter.id}"
}
resource "abiquo_fitpolicy" "fitpolicy" {
  policy     = "PERFORMANCE"
  target     = "${data.abiquo_datacenter.datacenter.id}"
}

resource "abiquo_computeload" "computeload" {
  cpuload    = "1000"
  ramload    = "95"
  target     = "${data.abiquo_datacenter.datacenter.id}"
}

resource "abiquo_storageload" "storageload" {
  load   = "95"
  target = "${data.abiquo_datacenter.datacenter.id}"
}
