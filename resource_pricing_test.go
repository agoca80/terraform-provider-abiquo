package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var pricingTestHelper = &testHelper{
	kind:  "abiquo_pricing",
	media: "pricingtemplate",
	config: `
	data "abiquo_datacenter" "test" {
		name = "datacenter 1"
	}

	data "abiquo_dstier" "test" {
		datacenter = "${data.abiquo_datacenter.test.id}"
		name       = "Default Tier"
	}

	data "abiquo_tier" "test" {
		datacenter = "${data.abiquo_datacenter.test.id}"
		name       = "Default Tier 1"
	}

	resource "abiquo_currency" "test" {
		digits = 2
		symbol = "TEST"
		name   = "TestAccPricing"
	}
	
	resource "abiquo_costcode" "test" {
		currency { href = "${abiquo_currency.test.id}", price = 1 }
		description = "TestAccPricing"
		name        = "TestAccPricing"
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

		costcode { href  = "${abiquo_costcode.test.id}", price = 7.9 }

		datacenter {
			href = "${data.abiquo_datacenter.test.id}"
			datastore_tier { href  = "${data.abiquo_dstier.test.id}", price = 2.3 }
			tier           { href  = "${data.abiquo_tier.test.id}"  , price = 4.5 }
			firewall = 1.2
		}
	}
	`,
}

func TestAccPricing_update(t *testing.T) {
	resource.Test(t, pricingTestHelper.updateCase(t))
}
