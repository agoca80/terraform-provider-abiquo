resource "abiquo_enterprise" "test" {
    name            = "test enterprise"
    pricingtemplate = "${abiquo_pricingtemplate.test.id}"
    properties = {
      "property0" = "value0"
      "property1" = "value1"
    }
}

data "abiquo_datacenter" "test" { name = "${var.datacenter}" }

resource "abiquo_currency" "test" {
  digits = 2
  symbol = "TEST"
  name   = "test enterprise"
}

resource "abiquo_pricingtemplate" "test" {
  currency               = "${abiquo_currency.test.id}"
  charging_period        = "DAY"
  deploy_message         = "test enterprise"
  description            = "test enterprise"
  minimum_charge         = 1
  minimum_charge_period  = "DAY"
  name                   = "test enterprise"
  standing_charge_period = 1
}
