# Provider configuration: Get these from the environment
variable "endpoint" { }
variable "username" { }
variable "password" { }

provider "abiquo" {
  endpoint       = "${var.endpoint}"
  username       = "${var.username}"
  password       = "${var.password}"
}

data "abiquo_datacenter" "dc1" { name = "datacenter 1" }
data "abiquo_datacenter" "dc2" { name = "datacenter 2" }

data "abiquo_dstier" "dc1dstier1" {
  datacenter = "${data.abiquo_datacenter.dc1.id}"
  name       = "Default Tier"
}

data "abiquo_dstier" "dc2dstier1" {
  datacenter = "${data.abiquo_datacenter.dc2.id}"
  name       = "Default Tier"
}

data "abiquo_tier" "tier1" {
  datacenter = "${data.abiquo_datacenter.dc1.id}"
  name       = "Default Tier 1"
}

data "abiquo_tier" "tier2" {
  datacenter = "${data.abiquo_datacenter.dc1.id}"
  name       = "Default Tier 2"
}

data "abiquo_tier" "tier3" {
  datacenter = "${data.abiquo_datacenter.dc1.id}"
  name       = "Default Tier 3"
}

data "abiquo_tier" "tier4" {
  datacenter = "${data.abiquo_datacenter.dc1.id}"
  name       = "Default Tier 4"
}

resource "abiquo_currency" "currency" {
  count  = 2
  digits = 0
  symbol = "TEST - T${count.index}"
  name   = "currency${count.index}"
}

resource "abiquo_costcode" "costcode" {
  currency { href = "${abiquo_currency.currency.0.id}", price = 1.2 }
  currency { href = "${abiquo_currency.currency.1.id}", price = 2.3 }
  description = "tf costcode${count.index}"
  name        = "tf costcode${count.index}"
  count       = 2
}

resource "abiquo_pricing" "pricing" {
  currency               = "${abiquo_currency.currency.0.id}"
  charging_period        = "DAY"
  description            = "pricing"
  minimum_charge         = 1
  minimum_charge_period  = "DAY"
  name                   = "pricing"
  standing_charge_period = 1

 costcode { href  = "${abiquo_costcode.costcode.0.id}", price = 7.9 }
 costcode { href  = "${abiquo_costcode.costcode.1.id}", price = 5.8 }

 datacenter {
   href = "${data.abiquo_datacenter.dc1.id}"
   datastore_tier { href  = "${data.abiquo_dstier.dc1dstier1.id}", price = 2.3 }
   tier           { href  = "${data.abiquo_tier.tier2.id}" , price = 8.9 }
   tier           { href  = "${data.abiquo_tier.tier3.id}" , price = 8.9 }
   tier           { href  = "${data.abiquo_tier.tier4.id}" , price = 1.2 }
   firewall = 2.3
 }

 datacenter {
   href = "${data.abiquo_datacenter.dc2.id}"
   datastore_tier { href  = "${data.abiquo_dstier.dc2dstier1.id}", price = 2.3 }
   firewall = 1.2
 }
}
