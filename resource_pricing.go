package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

var pricingDatacenterSchema = map[string]*schema.Schema{
	"href": &schema.Schema{
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
	"datastore_tier": &schema.Schema{
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"href": &schema.Schema{
					Required:     true,
					Type:         schema.TypeString,
					ValidateFunc: validateURL,
				},
				"price": &schema.Schema{
					Default:      0.0,
					Optional:     true,
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
		Default:  0.0,
		Optional: true,
		Type:     schema.TypeFloat,
	},
	"hd_gb": &schema.Schema{
		Default:  0.0,
		Optional: true,
		Type:     schema.TypeFloat,
	},
	"hardware_profile": &schema.Schema{
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"href": &schema.Schema{
					Required:     true,
					Type:         schema.TypeString,
					ValidateFunc: validateURL,
				},
				"price": &schema.Schema{
					Default:      0.0,
					Optional:     true,
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
		Default:  0.0,
		Optional: true,
		Type:     schema.TypeFloat,
	},
	"loadbalancer": &schema.Schema{
		Default:  0.0,
		Optional: true,
		Type:     schema.TypeFloat,
	},
	"memory": &schema.Schema{
		Default:  0.0,
		Optional: true,
		Type:     schema.TypeFloat,
	},
	"memory_on": &schema.Schema{
		Default:  0.0,
		Optional: true,
		Type:     schema.TypeFloat,
	},
	"memory_off": &schema.Schema{
		Default:  0.0,
		Optional: true,
		Type:     schema.TypeFloat,
	},
	"nat_ip": &schema.Schema{
		Default:  0.0,
		Optional: true,
		Type:     schema.TypeFloat,
	},
	"public_ip": &schema.Schema{
		Default:  0.0,
		Optional: true,
		Type:     schema.TypeFloat,
	},
	"repository": &schema.Schema{
		Default:  0.0,
		Optional: true,
		Type:     schema.TypeFloat,
	},
	"tier": &schema.Schema{
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"href": &schema.Schema{
					Required:     true,
					Type:         schema.TypeString,
					ValidateFunc: validateURL,
				},
				"price": &schema.Schema{
					Default:      0.0,
					Optional:     true,
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
		Default:  0.0,
		Optional: true,
		Type:     schema.TypeFloat,
	},
	"vcpu_on": &schema.Schema{
		Default:  0.0,
		Optional: true,
		Type:     schema.TypeFloat,
	},
	"vcpu_off": &schema.Schema{
		Default:  0.0,
		Optional: true,
		Type:     schema.TypeFloat,
	},
	"vlan": &schema.Schema{
		Default:  0.0,
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
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"href": &schema.Schema{
					Required:     true,
					Type:         schema.TypeString,
					ValidateFunc: validateURL,
				},
				"price": &schema.Schema{
					Default:      0.0,
					Optional:     true,
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
		Computed: true,
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
	if r == nil {
		return
	}
	resources := r.(*schema.Set)
	for _, r := range resources.List() {
		resource := r.(map[string]interface{})
		href := resource["href"].(string)
		rp = append(rp, abiquo.ResourcePrice{
			Price: resource["price"].(float64),
			DTO:   core.NewDTO(core.NewLinkType(href, media).SetRel(media)),
		})
	}
	return
}

func datacenterPrices(dc interface{}) abiquo.AbstractDCPrice {
	datacenter := dc.(map[string]interface{})
	href := datacenter["href"].(string)
	return abiquo.AbstractDCPrice{
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
	datacentersPrices := []abiquo.AbstractDCPrice{}
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
		MinimumCharge:       d.int("minimum_charge"),
		MinimumChargePeriod: pricingPeriod[d.string("minimum_charge_period")],
		DTO: core.NewDTO(
			d.linkTypeRel("currency", "currency", "currency"),
		),
	}
}

func pricingEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType("config/pricingtemplates", "pricingtemplate")
}

func resourcePricesRead(resources []abiquo.ResourcePrice, rel string) (set *schema.Set) {
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
