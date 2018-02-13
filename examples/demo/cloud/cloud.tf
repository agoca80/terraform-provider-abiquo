# Provider configuration: Get these from the environment
# Terraform environment variables look like TF_VAR_name, so, to configure the
# properties below, you need to configure the following varables:
# TF_VAR_endpoint
# TF_VAR_username
# TF_VAR_password
# Should these environment variables not exist, terraform will ask their values
# during the initialization.
variable "endpoint" { }
variable "username" { }
variable "password" { }

provider "abiquo" {
  endpoint       = "${var.endpoint}"
  username       = "${var.username}"
  password       = "${var.password}"
}

data "abiquo_datacenter" "datacenter"  { name = "datacenter 1" }

resource "abiquo_enterprise" "enterprise" {
   name = "terraform demo"
   cpusoft  = 6    , cpuhard  = 8
   ramsoft  = 8192 , ramhard  = 16384
   vlansoft = 2    , vlanhard = 4
}

# PENDING Limits should be a set of resources inside enterprise
resource "abiquo_limit" "limit" {
  enterprise = "${abiquo_enterprise.enterprise.id}"
  location   = "${data.abiquo_datacenter.datacenter.id}"
}
