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
	"name":                attribute(required, text),
	"algorithm":           attribute(required, label(algorithms)),
	"internal":            attribute(optional, boolean),
	"routingrules":        attribute(required, list(lbRuleResource), min(1)),
	"privatenetwork":      attribute(optional, forceNew, link("privatenetwork")),
	"virtualdatacenter":   attribute(required, forceNew, link("virtualdatacenter")),
	"loadbalanceraddress": attribute(computed, text),
	"virtualmachines":     attribute(computed, list(attribute(text))),
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
		Name:      d.string("name"),
		Algorithm: d.string("algorithm"),
		LoadBalancerAddresses: abiquo.LoadBalancerAddresses{
			Collection: []abiquo.LoadBalancerAddress{
				abiquo.LoadBalancerAddress{Internal: d.boolean("internal")},
			},
		},
		RoutingRules: abiquo.RoutingRules{
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

func lbUpdate(rd *schema.ResourceData, m interface{}) (err error) {
	d := newResourceData(rd, "loadbalancer")
	lb := lbNew(d).(*abiquo.LoadBalancer)
	if err = core.Update(d, lb); err == nil {
		err = lb.SetRules(lbRules(d))
	}
	return
}

func lbEndpoint(d *resourceData) (link *core.Link) {
	if device := vdcDevice(d.link("virtualdatacenter")); device != nil {
		link = device.Rel("loadbalancers").SetType("loadbalancer")
	}
	return
}
