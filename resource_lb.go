package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var algorithms = []string{"Default", "ROUND_ROBIN", "LEAST_CONNECTIONS", "SOURCE_IP"}

var lbRuleResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"protocolin":  attribute(required, protocol),
		"protocolout": attribute(required, protocol),
		"portout":     attribute(required, port),
		"portin":      attribute(required, port),
	},
}

var lbSchema = map[string]*schema.Schema{
	"virtualdatacenter":   endpoint("virtualdatacenter"),
	"name":                attribute(required, text),
	"algorithm":           attribute(required, label(algorithms)),
	"internal":            attribute(optional, boolean),
	"routingrules":        attribute(required, list(lbRuleResource), min(1)),
	"privatenetwork":      attribute(optional, forceNew, link("privatenetwork")),
	"loadbalanceraddress": attribute(computed, text),
	"virtualmachines":     attribute(computed, list(text)),
}

func lbRules(d *resourceData) (rules []abiquo.LoadBalancerRule) {
	for _, r := range d.slice("routingrules") {
		rule := abiquo.LoadBalancerRule{}
		mapDecoder(r, &rule)
		rules = append(rules, rule)
	}
	return
}

func lbNew(d *resourceData) core.Resource {
	return &abiquo.LoadBalancer{
		Name:      d.string("name"),
		Algorithm: d.string("algorithm"),
		LoadBalancerAddresses: abiquo.LoadBalancerAddresses{
			Collection: []abiquo.LoadBalancerAddress{
				abiquo.LoadBalancerAddress{Internal: d.boolean("internal")},
			},
		},
		RoutingRules: abiquo.LoadBalancerRules{
			Collection: lbRules(d),
		},
		DTO: core.NewDTO(
			d.link("virtualdatacenter"),
			d.link("privatenetwork").SetType("vlan"),
		),
	}
}

func lbRead(d *resourceData, resource core.Resource) (err error) {
	lb := resource.(*abiquo.LoadBalancer)

	// Get lb virtualmachines hrefs
	virtualmachines := []interface{}{}
	lb.VMs().Map(func(l *core.Link) {
		virtualmachines = append(virtualmachines, l.Href)
	})

	d.Set("name", lb.Name)
	d.Set("algorithm", lb.Algorithm)
	d.Set("loadbalanceraddress", lb.Rel("loadbalanceraddress").Title)
	d.Set("virtualmachines", virtualmachines)

	return
}

func lbUpdate(d *resourceData, resource core.Resource) (err error) {
	if d.HasChange("routingrules") {
		loadBalancerRules := abiquo.LoadBalancerRules{
			Collection: lbRules(d),
		}
		if err = core.Update(resource.Rel("rules"), loadBalancerRules); err != nil {
			return
		}
	}
	return
}

func lbEndpoint(d *resourceData) (link *core.Link) {
	if device := vdcDevice(d.link("virtualdatacenter")); device != nil {
		link = device.Rel("loadbalancers").SetType("loadbalancer")
	}
	return
}
