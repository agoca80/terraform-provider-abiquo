# Provider configuration: Get these from the environment
variable "endpoint" { }

provider "abiquo" {
  endpoint       = "${var.endpoint}"
  username       = "hwprofiles"
  password       = "12qwaszx"
}

