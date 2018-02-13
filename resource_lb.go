package main

import (
	"github.com/abiquo/opal/abiquo"
	"github.com/abiquo/opal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var algorithms = []string{"ROUND_ROBIN", "LEAST_CONNECTIONS", "SOURCE_IP"}

var lbResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name":      Required().String(),
		"algorithm": Required().ValidateString(algorithms),
		"routingrules": required(&schema.Schema{Type: schema.TypeList, MinItems: 1, Elem: &schema.Resource{
			Schema: lbRuleSchema,
		}}),
		"privatenetwork":    Required().Renew().Link(),
		"virtualdatacenter": Required().Renew().Link(),
	},
	Delete: resourceDelete,
	Exists: resourceExists("loadbalancer"),
	Create: resourceCreate(lbNew, lbEndpoint),
	Update: resourceUpdate(lbNew, "loadbalancer"),
	Read:   resourceRead(lbNew, lbRead, "loadbalancer"),
}

var lbRuleSchema = map[string]*schema.Schema{
	"protocolin":  required(validate(fieldString(), validateString([]string{"TCP", "HTTP", "HTTPS"}))),
	"protocolout": required(validate(fieldString(), validateString([]string{"TCP", "HTTP", "HTTPS"}))),
	"portout":     Required().Number(),
	"portin":      Required().Number(),
}

func lbAddresses(d *resourceData) abiquo.LoadBalancerAddresses {
	return abiquo.LoadBalancerAddresses{
		Collection: []abiquo.LoadBalancerAddress{
			abiquo.LoadBalancerAddress{Internal: false},
		},
	}
}

func lbRules(d *resourceData) (rules []abiquo.RoutingRule) {
	for _, r := range d.slice("routingrules") {
		rule := abiquo.RoutingRule{}
		mapDecoder(r, &rule)
		rules = append(rules, rule)
	}
	return
}

func lbNew(d *resourceData) core.Resource {
	return &abiquo.LoadBalancer{
		Name:                  d.string("name"),
		Algorithm:             d.string("algorithm"),
		LoadBalancerAddresses: lbAddresses(d),
		RoutingRules: abiquo.RoutingRules{
			Collection: lbRules(d),
		},
		DTO: core.NewDTO(
			d.link("virtualdatacenter"),
			d.link("privatenetwork"),
		),
	}
}

func lbRead(d *resourceData, resource core.Resource) (err error) {
	lb := resource.(*abiquo.LoadBalancer)
	d.Set("name", lb.Name)
	d.Set("algorithm", lb.Algorithm)
	return
}

func lbUpdate(rd *schema.ResourceData, m interface{}) (err error) {
	d := newResourceData(rd, "loadbalancer")
	lb := lbNew(d).(*abiquo.LoadBalancer)
	if err = core.Update(d, lb); err == nil {
		err = lb.SetRules(lbRules(d))
	}
	return
}

func lbEndpoint(d *resourceData) (link *core.Link) {
	vdc := new(abiquo.VirtualDatacenter)
	if core.Read(d.link("virtualdatacenter"), vdc) == nil {
		endpoint := vdc.Rel("device")
		if endpoint == nil {
			return nil
		}

		device := new(abiquo.Device)
		if core.Read(endpoint, device) == nil {
			link = device.Rel("loadbalancers").SetType("loadbalancer")
		}
	}
	return
}
