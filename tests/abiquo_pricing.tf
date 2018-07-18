resource "abiquo_pricing" "test" {
  currency               = "${abiquo_currency.test.id}"
  charging_period        = "DAY"
  deploy_message         = "test pricing"
  description            = "test pricing"
  minimum_charge         = 1
  minimum_charge_period  = "DAY"
  name                   = "test pricing"
  standing_charge_period = 1

  costcode {
    href  = "${abiquo_costcode.test.id}"
    price = 7.9
  }

  datacenter {
    href = "${data.abiquo_datacenter.test.id}"
    datastore_tier { href  = "${data.abiquo_dstier.test.id}", price = 2.3 }
    tier           { href  = "${data.abiquo_tier.test.id}",   price = 4.5 }
    firewall = 1.2
  }
}

data "abiquo_datacenter" "test" {
  name = "datacenter 1"
}

data "abiquo_dstier" "test" {
  datacenter = "${data.abiquo_datacenter.test.id}"
  name       = "Default Tier"
}

data "abiquo_tier" "test" {
  location = "${data.abiquo_datacenter.test.tiers}"
  name     = "Default Tier 1"
}

resource "abiquo_currency" "test" {
  digits = 2
  symbol = "TEST"
  name   = "test pricing"
}

resource "abiquo_costcode" "test" {
  currency {
    href = "${abiquo_currency.test.id}"

    price = 1
  }

  description = "test pricing"
  name        = "test pricing"
}
