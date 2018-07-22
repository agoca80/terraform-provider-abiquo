package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var currencySchema = map[string]*schema.Schema{
	"digits": attribute(required, natural),
	"name":   attribute(required, text),
	"symbol": attribute(required, text),
}

func currencyNew(d *resourceData) core.Resource {
	return &abiquo.Currency{
		Digits: d.integer("digits"),
		Name:   d.string("name"),
		Symbol: d.string("symbol"),
	}
}

func currencyRead(d *resourceData, resource core.Resource) (err error) {
	currency := resource.(*abiquo.Currency)
	d.Set("digits", currency.Digits)
	d.Set("name", currency.Name)
	d.Set("symbol", currency.Symbol)
	return
}

var currency = &description{
	dto:      currencyNew,
	endpoint: endpointConst("config/currencies"),
	media:    "currency",
	read:     currencyRead,
	Resource: &schema.Resource{Schema: currencySchema},
}
