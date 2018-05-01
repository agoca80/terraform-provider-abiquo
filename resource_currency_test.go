package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var currencyTestHelper = &testHelper{
	kind:  "abiquo_currency",
	media: "currency",
	config: `
	resource "abiquo_currency" "test" {
		digits = 2
		symbol = "TEST"
	  name   = "testAccCurrencyBasic"
	}
	`,
}

func TestAccCurrency_update(t *testing.T) {
	resource.Test(t, currencyTestHelper.updateCase(t))
}
