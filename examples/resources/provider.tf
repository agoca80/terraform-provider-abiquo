# Provider configuration: Get these from the environment
variable "endpoint" { }
variable "username" { }
variable "password" { }

provider "abiquo" {
  endpoint       = "${var.endpoint}"
  username       = "${var.username}"
  password       = "${var.password}"
}

data "abiquo_enterprise" "enterprise" { name = "Abiquo" }
data "abiquo_datacenter" "datacenter" { name = "datacenter 1" }
data "abiquo_location"   "location"   { name = "datacenter 1" }
