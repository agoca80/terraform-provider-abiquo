package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var costCodeSchema = map[string]*schema.Schema{
	"currency":    attribute(required, prices, min(1)),
	"name":        attribute(required, text),
	"description": attribute(optional, text),
}

func costCodeNew(d *resourceData) core.Resource {
	currencies := []abiquo.PricingResource{}
	for _, c := range d.set("currency").List() {
		currency := c.(map[string]interface{})
		href := currency["href"].(string)
		link := core.NewLinkType(href, "currency").SetRel("currency")
		currencies = append(currencies, abiquo.PricingResource{
			Price: currency["price"].(float64),
			DTO:   core.NewDTO(link),
		})
	}
	return &abiquo.CostCode{
		CurrencyPrices: currencies,
		Name:           d.string("name"),
		Description:    d.string("description"),
	}
}

func costCodeRead(d *resourceData, resource core.Resource) (err error) {
	costCode := resource.(*abiquo.CostCode)
	currency := []interface{}{}
	for _, c := range costCode.CurrencyPrices {
		currency = append(currency, map[string]interface{}{
			"price":    c.Price,
			"currency": c.URL(),
		})
	}
	d.Set("description", costCode.Description)
	d.Set("name", costCode.Name)
	d.Set("currency", currency)
	return
}

var costcode = &description{
	Resource: &schema.Resource{Schema: costCodeSchema},
	dto:      costCodeNew,
	endpoint: endpointConst("config/costcodes"),
	media:    "costcode",
	read:     costCodeRead,
}
