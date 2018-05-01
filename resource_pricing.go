package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

var pricingDatacenterSchema = map[string]*schema.Schema{
	"href": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
	"datastore_tier": &schema.Schema{
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"href": &schema.Schema{
					Required:     true,
					Type:         schema.TypeString,
					ValidateFunc: validateURL,
				},
				"price": &schema.Schema{
					Required:     true,
					Type:         schema.TypeFloat,
					ValidateFunc: validatePrice,
				},
			},
		},
		Optional: true,
		Set:      resourceSet,
		Type:     schema.TypeSet,
	},
	"firewall": &schema.Schema{
		Optional: true,
		Type:     schema.TypeFloat,
	},
	"hd_gb": &schema.Schema{
		Optional: true,
		Type:     schema.TypeFloat,
	},
	"hardware_profile": &schema.Schema{
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"href": &schema.Schema{
					Required:     true,
					Type:         schema.TypeString,
					ValidateFunc: validateURL,
				},
				"price": &schema.Schema{
					Required:     true,
					Type:         schema.TypeFloat,
					ValidateFunc: validatePrice,
				},
			},
		},
		Optional: true,
		Set:      resourceSet,
		Type:     schema.TypeSet,
	},
	"layer": &schema.Schema{
		Optional: true,
		Type:     schema.TypeFloat,
	},
	"loadbalancer": &schema.Schema{
		Optional: true,
		Type:     schema.TypeFloat,
	},
	"memory": &schema.Schema{
		Optional: true,
		Type:     schema.TypeFloat,
	},
	"memory_on": &schema.Schema{
		Optional: true,
		Type:     schema.TypeFloat,
	},
	"memory_off": &schema.Schema{
		Optional: true,
		Type:     schema.TypeFloat,
	},
	"nat_ip": &schema.Schema{
		Optional: true,
		Type:     schema.TypeFloat,
	},
	"public_ip": &schema.Schema{
		Optional: true,
		Type:     schema.TypeFloat,
	},
	"repository": &schema.Schema{
		Optional: true,
		Type:     schema.TypeFloat,
	},
	"tier": &schema.Schema{
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"href": &schema.Schema{
					Required:     true,
					Type:         schema.TypeString,
					ValidateFunc: validateURL,
				},
				"price": &schema.Schema{
					Required:     true,
					Type:         schema.TypeFloat,
					ValidateFunc: validatePrice,
				},
			},
		},
		Optional: true,
		Set:      resourceSet,
		Type:     schema.TypeSet,
	},
	"vcpu": &schema.Schema{
		Optional: true,
		Type:     schema.TypeFloat,
	},
	"vcpu_on": &schema.Schema{
		Optional: true,
		Type:     schema.TypeFloat,
	},
	"vcpu_off": &schema.Schema{
		Optional: true,
		Type:     schema.TypeFloat,
	},
	"vlan": &schema.Schema{
		Optional: true,
		Type:     schema.TypeFloat,
	},
}

var pricingSchema = map[string]*schema.Schema{
	"charging_period": &schema.Schema{
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validation.StringInSlice([]string{"DAY", "WEEK", "MONTH", "QUARTER", "YEAR"}, false),
	},
	"costcode": &schema.Schema{
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"href": &schema.Schema{
					Required:     true,
					Type:         schema.TypeString,
					ValidateFunc: validateURL,
				},
				"price": &schema.Schema{
					Required:     true,
					Type:         schema.TypeFloat,
					ValidateFunc: validatePrice,
				},
			},
		},
		Optional: true,
		Set:      resourceSet,
		Type:     schema.TypeSet,
	},
	"currency": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
	"datacenter": &schema.Schema{
		Elem: &schema.Resource{
			Schema: pricingDatacenterSchema,
		},
		Optional: true,
		Set:      resourceSet,
		Type:     schema.TypeSet,
	},
	"deploy_message": &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	},
	"description": &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	},
	"minimum_charge": &schema.Schema{
		Required: true,
		Type:     schema.TypeInt,
	},
	"minimum_charge_period": &schema.Schema{
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validation.StringInSlice([]string{"MINUTE", "HOUR", "DAY", "WEEK", "MONTH", "QUARTER", "YEAR"}, false),
	},
	"name": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"show_charges_before": &schema.Schema{
		Optional: true,
		Type:     schema.TypeBool,
	},
	"show_minimun_charge": &schema.Schema{
		Optional: true,
		Type:     schema.TypeBool,
	},
	"standing_charge_period": &schema.Schema{
		Optional: true,
		Type:     schema.TypeInt,
	},
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

func resourcePrices(r interface{}, media string) (rp []abiquo.ResourcePrice) {
	resources := r.(*schema.Set)
	for _, r := range resources.List() {
		resource := r.(map[string]interface{})
		href := resource["href"].(string)
		link := core.NewLinkType(href, media).SetRel(media)
		rp = append(rp, abiquo.ResourcePrice{
			Price: resource["price"].(float64),
			DTO:   core.NewDTO(link),
		})
	}
	return
}

func pricingNew(d *resourceData) core.Resource {
	currency := d.linkTypeRel("currency", "currency", "currency")
	datacenters := []abiquo.AbstractDCPrice{}
	for _, dc := range d.set("datacenter").List() {
		datacenter := dc.(map[string]interface{})
		href := datacenter["href"].(string)
		link := core.NewLinkType(href, "datacenter").SetRel("datacenter")
		datacenters = append(datacenters, abiquo.AbstractDCPrice{
			Firewall:       datacenter["firewall"].(float64),
			HardDiskGB:     datacenter["hd_gb"].(float64),
			HPAbstractDC:   resourcePrices(datacenter["hardware_profile"], "hardwareprofile"),
			DatastoreTiers: resourcePrices(datacenter["datastore_tier"], "datastoretier"),
			Layer:          datacenter["layer"].(float64),
			LoadBalancer:   datacenter["loadbalancer"].(float64),
			MemoryGB:       datacenter["memory"].(float64),
			MemoryOnGB:     datacenter["memory_on"].(float64),
			MemoryOffGB:    datacenter["memory_off"].(float64),
			NatIP:          datacenter["nat_ip"].(float64),
			PublicIP:       datacenter["public_ip"].(float64),
			RepositoryGB:   datacenter["repository"].(float64),
			Tiers:          resourcePrices(datacenter["tier"], "tier"),
			VCPU:           datacenter["vcpu"].(float64),
			VCPUOn:         datacenter["vcpu_on"].(float64),
			VCPUOff:        datacenter["vcpu_off"].(float64),
			VLAN:           datacenter["vlan"].(float64),
			DTO:            core.NewDTO(link),
		})
	}

	costCodes := []abiquo.ResourcePrice{}
	for _, c := range d.set("costcode").List() {
		costCode := c.(map[string]interface{})
		href := costCode["href"].(string)
		link := core.NewLinkType(href, "costcode").SetRel("costcode")
		costCodes = append(costCodes, abiquo.ResourcePrice{
			Price: costCode["price"].(float64),
			DTO:   core.NewDTO(link, currency),
		})
	}

	return &abiquo.PricingTemplate{
		AbstractDCPrices:    datacenters,
		ChargingPeriod:      pricingPeriod[d.string("charging_period")],
		CostCodes:           costCodes,
		Name:                d.string("name"),
		Description:         d.string("description"),
		MinimumCharge:       d.int("minimum_charge"),
		MinimumChargePeriod: pricingPeriod[d.string("minimum_charge_period")],
		DTO:                 core.NewDTO(currency),
	}
}

func pricingEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType("config/pricingtemplates", "pricingtemplate")
}

func resourcePricesRead(resources []abiquo.ResourcePrice) (prices []interface{}) {
	for _, resource := range resources {
		prices = append(prices, map[string]interface{}{
			"href":  resource.URL(),
			"price": resource.Price,
		})
	}
	return
}

func pricingRead(d *resourceData, resource core.Resource) (err error) {
	pricing := resource.(*abiquo.PricingTemplate)
	datacenters := []interface{}{}
	for _, d := range pricing.AbstractDCPrices {
		datacenters = append(datacenters, map[string]interface{}{
			"datastore_tier":   resourcePricesRead(d.DatastoreTiers),
			"firewall":         d.Firewall,
			"hardware_profile": resourcePricesRead(d.HPAbstractDC),
			"hdgb":             d.HardDiskGB,
			"layer":            d.Layer,
			"loadbalancer":     d.LoadBalancer,
			"memory":           d.MemoryGB,
			"memoryOn":         d.MemoryOnGB,
			"memoryOff":        d.MemoryOffGB,
			"natip":            d.NatIP,
			"publicip":         d.PublicIP,
			"repository":       d.RepositoryGB,
			"tier":             resourcePricesRead(d.Tiers),
			"vcpu":             d.VCPU,
			"vcpuon":           d.VCPUOn,
			"vcpuoff":          d.VCPUOff,
			"vlan":             d.VLAN,
			"datacenter":       d.URL(),
		})
	}
	d.Set("description", pricing.Description)
	d.Set("costcode", resourcePricesRead(pricing.CostCodes))
	d.Set("name", pricing.Name)
	return
}
