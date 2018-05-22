package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var enterpriseTestHelper = &testHelper{
	kind:  "abiquo_enterprise",
	media: "enterprise",
	config: `
	data "abiquo_datacenter" "test" { name = "datacenter 1" }

	resource "abiquo_currency" "test" {
		digits = 2
		symbol = "TEST"
		name   = "TestAccEnterprise"
	}

	resource "abiquo_pricing" "test" {
		currency               = "${abiquo_currency.test.id}"
		charging_period        = "DAY"
		deploy_message         = "TestAccPricing"
		description            = "TestAccPricing"
		minimum_charge         = 1
		minimum_charge_period  = "DAY"
		name                   = "TestAccPricing"
		standing_charge_period = 1
	}

	resource "abiquo_enterprise" "test" {
	    name            = "TestAccEnterprise"
			pricingtemplate = "${abiquo_pricing.test.id}"
			properties = {
				"property0" = "value0"
				"property1" = "value1"
			}
	}
	`,
}

func TestAccEnterprise_update(t *testing.T) {
	resource.Test(t, enterpriseTestHelper.updateCase(t))
}
