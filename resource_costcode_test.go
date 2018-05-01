package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var costCodeTestHelper = &testHelper{
	kind:  "abiquo_costcode",
	media: "costcode",
	config: `
	resource "abiquo_currency" "a" {
		digits = 1
		symbol = "TEST - A"
		name   = "testAccCostCode - A"
	}

	resource "abiquo_currency" "b" {
		digits = 2
		symbol = "TEST - B"
		name   = "testAccCostCode - B"
	}

	resource "abiquo_costcode" "test" {
		currency { href = "${abiquo_currency.a.id}", price = 1 }
		currency { href = "${abiquo_currency.b.id}", price = 2 }
		description = "testAccCostCode"
	  name        = "testAccCostCode"
	}
	`,
}

func TestAccCostCode_update(t *testing.T) {
	resource.Test(t, costCodeTestHelper.updateCase(t))
}
