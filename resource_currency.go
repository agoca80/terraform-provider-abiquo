package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var currencySchema = map[string]*schema.Schema{
	"digits": &schema.Schema{
		Required: true,
		Type:     schema.TypeInt,
		ValidateFunc: func(d interface{}, key string) (strs []string, errs []error) {
			if 0 > d.(int) || 2 < d.(int) {
				errs = append(errs, fmt.Errorf("digits should be between 0 and 2"))
			}
			return
		},
	},
	"name": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"symbol": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
}

func currencyNew(d *resourceData) core.Resource {
	return &abiquo.Currency{
		Digits: d.int("digits"),
		Name:   d.string("name"),
		Symbol: d.string("symbol"),
	}
}

func currencyEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType("config/currencies", "currency")
}

func currencyRead(d *resourceData, resource core.Resource) (err error) {
	currency := resource.(*abiquo.Currency)
	d.Set("digits", currency.Digits)
	d.Set("name", currency.Name)
	d.Set("symbol", currency.Symbol)
	return
}
