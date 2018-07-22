package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var lbAlgorithms = []string{"Default", "ROUND_ROBIN", "LEAST_CONNECTIONS", "SOURCE_IP"}

var lbSchema = map[string]*schema.Schema{
	"device":              endpoint("device"),
	"name":                attribute(required, text),
	"algorithm":           attribute(required, label(lbAlgorithms)),
	"internal":            attribute(optional, boolean),
	"privatenetwork":      attribute(optional, link("privatenetwork"), forceNew),
	"loadbalanceraddress": attribute(computed, text),
	"virtualmachines":     attribute(computed, list(text)),
	"routingrules": attribute(required, min(1), list(&schema.Resource{
		Schema: map[string]*schema.Schema{
			"protocolin":  attribute(required, protocol),
			"protocolout": attribute(required, protocol),
			"portout":     attribute(required, port),
			"portin":      attribute(required, port),
		},
	})),
}

func lbRules(d *resourceData) (rules *abiquo.LoadBalancerRules) {
	rules = new(abiquo.LoadBalancerRules)
	for _, r := range d.slice("routingrules") {
		rule := abiquo.LoadBalancerRule{}
		mapDecoder(r, &rule)
		rules.Collection = append(rules.Collection, rule)
	}
	return
}

func lbNew(d *resourceData) core.Resource {
	return &abiquo.LoadBalancer{
		Name:      d.string("name"),
		Algorithm: d.string("algorithm"),
		Addresses: &abiquo.LoadBalancerAddresses{
			Collection: []abiquo.LoadBalancerAddress{
				abiquo.LoadBalancerAddress{Internal: d.boolean("internal")},
			},
		},
		Rules: lbRules(d),
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
		err = core.Update(resource.Rel("rules"), lbRules(d))
		if err != nil {
			return
		}
	}
	return
}

var loadbalancer = &description{
	media:    "loadbalancer",
	dto:      lbNew,
	endpoint: endpointPath("device", "/loadbalancers"),
	read:     lbRead,
	update:   lbUpdate,
	Resource: &schema.Resource{Schema: lbSchema},
}
