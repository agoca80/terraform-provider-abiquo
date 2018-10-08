package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var pricingDCResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"href":             attribute(required, href),
		"datastore_tier":   attribute(optional, computed, prices),
		"firewall":         attribute(price),
		"hd_gb":            attribute(price),
		"hardware_profile": attribute(optional, computed, prices),
		"layer":            attribute(price),
		"loadbalancer":     attribute(price),
		"memory":           attribute(price),
		"memory_on":        attribute(price),
		"memory_off":       attribute(price),
		"nat_ip":           attribute(price),
		"public_ip":        attribute(price),
		"repository":       attribute(price),
		"tier":             attribute(optional, computed, prices),
		"vcpu":             attribute(price),
		"vcpu_on":          attribute(price),
		"vcpu_off":         attribute(price),
		"vlan":             attribute(price),
	},
}

var pricingPeriodLabel = []string{"MINUTE", "HOUR", "DAY", "WEEK", "MONTH", "QUARTER", "YEAR"}

var pricingSchema = map[string]*schema.Schema{
	"charging_period":        attribute(required, label(pricingPeriodLabel[2:])),
	"costcode":               attribute(optional, computed, prices),
	"currency":               attribute(required, link("currency"), forceNew),
	"datacenter":             attribute(optional, computed, setFn(pricingDCResource, resourceSet)),
	"deploy_message":         attribute(optional, text),
	"description":            attribute(optional, text),
	"minimum_charge":         attribute(required, natural),
	"minimum_charge_period":  attribute(required, label(pricingPeriodLabel)),
	"name":                   attribute(required, text),
	"show_charges_before":    attribute(optional, boolean),
	"show_minimun_charge":    attribute(optional, boolean),
	"standing_charge_period": attribute(optional, integer),
}

var pricingPeriod = map[string]int{
	"MINUTE":  0,
	"HOUR":    1,
	"DAY":     2,
	"WEEK":    3,
	"MONTH":   4,
	"QUARTER": 5,
	"YEAR":    6,
}

func resourcePrices(r interface{}, media string) (rp []abiquo.PricingResource) {
	if r == nil {
		return
	}
	resources := r.(*schema.Set)
	for _, r := range resources.List() {
		resource := r.(map[string]interface{})
		href := resource["href"].(string)
		rp = append(rp, abiquo.PricingResource{
			Price: resource["price"].(float64),
			DTO:   core.NewDTO(core.NewLinkType(href, media).SetRel(media)),
		})
	}
	return
}

func datacenterPrices(dc interface{}) abiquo.PricingDatacenter {
	datacenter := dc.(map[string]interface{})
	href := datacenter["href"].(string)
	return abiquo.PricingDatacenter{
		HPAbstractDC:   resourcePrices(datacenter["hardware_profile"], "hardwareprofile"),
		DatastoreTiers: resourcePrices(datacenter["datastore_tier"], "datastoretier"),
		Tiers:          resourcePrices(datacenter["tier"], "tier"),
		Firewall:       datacenter["firewall"].(float64),
		HardDiskGB:     datacenter["hd_gb"].(float64),
		Layer:          datacenter["layer"].(float64),
		LoadBalancer:   datacenter["loadbalancer"].(float64),
		MemoryGB:       datacenter["memory"].(float64),
		MemoryOnGB:     datacenter["memory_on"].(float64),
		MemoryOffGB:    datacenter["memory_off"].(float64),
		NatIP:          datacenter["nat_ip"].(float64),
		PublicIP:       datacenter["public_ip"].(float64),
		RepositoryGB:   datacenter["repository"].(float64),
		VCPU:           datacenter["vcpu"].(float64),
		VCPUOn:         datacenter["vcpu_on"].(float64),
		VCPUOff:        datacenter["vcpu_off"].(float64),
		VLAN:           datacenter["vlan"].(float64),
		DTO:            core.NewDTO(core.NewLinkType(href, "datacenter").SetRel("datacenter")),
	}
}

func pricingNew(d *resourceData) core.Resource {
	datacentersPrices := []abiquo.PricingDatacenter{}
	if dcSet := d.set("datacenter"); dcSet != nil {
		for _, dc := range dcSet.List() {
			datacentersPrices = append(datacentersPrices, datacenterPrices(dc))
		}
	}

	return &abiquo.PricingTemplate{
		AbstractDCPrices:    datacentersPrices,
		ChargingPeriod:      pricingPeriod[d.string("charging_period")],
		CostCodes:           resourcePrices(d.Get("costcode"), "costcode"),
		Name:                d.string("name"),
		Description:         d.string("description"),
		MinimumCharge:       d.integer("minimum_charge"),
		MinimumChargePeriod: pricingPeriod[d.string("minimum_charge_period")],
		DTO: core.NewDTO(
			d.link("currency"),
		),
	}
}

func resourcePricesRead(resources []abiquo.PricingResource, rel string) (set *schema.Set) {
	set = schema.NewSet(resourceSet, nil)
	for _, resource := range resources {
		set.Add(map[string]interface{}{
			"href":  resource.Rel(rel).URL(),
			"price": resource.Price,
		})
	}
	return
}

func pricingRead(d *resourceData, resource core.Resource) (err error) {
	pricing := resource.(*abiquo.PricingTemplate)
	datacenters := schema.NewSet(resourceSet, nil)
	for _, dc := range pricing.AbstractDCPrices {
		datacenters.Add(map[string]interface{}{
			"tier":             resourcePricesRead(dc.Tiers, "tier"),
			"datastore_tier":   resourcePricesRead(dc.DatastoreTiers, "datastoretier"),
			"hardware_profile": resourcePricesRead(dc.HPAbstractDC, "hardwareprofile"),
			"firewall":         dc.Firewall,
			"hdgb":             dc.HardDiskGB,
			"layer":            dc.Layer,
			"loadbalancer":     dc.LoadBalancer,
			"memory":           dc.MemoryGB,
			"memoryOn":         dc.MemoryOnGB,
			"memoryOff":        dc.MemoryOffGB,
			"natip":            dc.NatIP,
			"publicip":         dc.PublicIP,
			"repository":       dc.RepositoryGB,
			"vcpu":             dc.VCPU,
			"vcpuon":           dc.VCPUOn,
			"vcpuoff":          dc.VCPUOff,
			"vlan":             dc.VLAN,
			"href":             dc.Rel("datacenter").URL(),
		})
	}
	d.Set("datacenter", datacenters)
	d.Set("description", pricing.Description)
	d.Set("costcode", resourcePricesRead(pricing.CostCodes, "costcode"))
	d.Set("name", pricing.Name)
	return
}

var pricingtemplate = &description{
	Resource: &schema.Resource{Schema: pricingSchema},
	dto:      pricingNew,
	endpoint: endpointConst("config/pricingtemplates"),
	media:    "pricingtemplate",
	read:     pricingRead,
}
