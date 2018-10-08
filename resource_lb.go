package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var lbAlgorithms = []string{"Default", "ROUND_ROBIN", "LEAST_CONNECTIONS", "SOURCE_IP"}

var lbSchema = map[string]*schema.Schema{
	"device":          endpoint("device"),
	"name":            attribute(required, text),
	"algorithm":       attribute(required, label(lbAlgorithms)),
	"internal":        attribute(optional, boolean),
	"external":        attribute(optional, boolean),
	"privatenetwork":  attribute(optional, link("privatenetwork"), forceNew),
	"externalips":     attribute(computed, list(text)),
	"internalips":     attribute(computed, list(text)),
	"virtualmachines": attribute(computed, list(text)),
	"routingrules": attribute(required, min(1), list(&schema.Resource{
		Schema: map[string]*schema.Schema{
			"protocolin":  attribute(required, label([]string{"TCP", "HTTP", "HTTPS"})),
			"protocolout": attribute(required, label([]string{"TCP", "HTTP", "HTTPS"})),
			"portout":     attribute(required, port),
			"portin":      attribute(required, port),
		},
	})),
	"healthchecks": attribute(optional, min(0), list(&schema.Resource{
		Schema: map[string]*schema.Schema{
			//"id":       attribute(computed, text),
			"name":         attribute(required, text),
			"attempts":     attribute(required, positive),
			"intervalinms": attribute(required, positive),
			"timeoutinms":  attribute(required, positive),
			"protocol":     attribute(required, label([]string{"TCP", "PING", "HTTP", "HTTPS"})),
			"path":         attribute(optional, text),
			"port":         attribute(required, port),
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

func lbHealthChecks(d *resourceData) (healthChecks *abiquo.LoadBalancerHealthChecks) {
	healthChecks = new(abiquo.LoadBalancerHealthChecks)
	for _, r := range d.slice("healthchecks") {
		healthCheck := abiquo.LoadBalancerHealthCheck{}
		mapDecoder(r, &healthCheck)
		healthChecks.Collection = append(healthChecks.Collection, healthCheck)
	}
	return
}

func lbAddresses(d *resourceData) (addresses *abiquo.LoadBalancerAddresses) {
	addresses = &abiquo.LoadBalancerAddresses{
		Collection: []abiquo.LoadBalancerAddress{},
	}
	if d.boolean("internal") {
		addresses.Collection = append(addresses.Collection, abiquo.LoadBalancerAddress{Internal: true})
	}
	if d.boolean("external") {
		addresses.Collection = append(addresses.Collection, abiquo.LoadBalancerAddress{Internal: false})
	}
	return
}

func lbNew(d *resourceData) core.Resource {
	return &abiquo.LoadBalancer{
		Addresses:    lbAddresses(d),
		Algorithm:    d.string("algorithm"),
		Name:         d.string("name"),
		HealthChecks: lbHealthChecks(d),
		Rules:        lbRules(d),
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
	d.Set("virtualmachines", virtualmachines)

	addresses := &abiquo.LoadBalancerAddresses{}
	if err = core.Read(lb.Rel("addresses"), addresses); err != nil {
		return
	}
	d.Set("externalips", addresses.Endpoints(false))
	d.Set("internalips", addresses.Endpoints(true))
	return
}

func lbUpdate(d *resourceData, resource core.Resource) (err error) {
	if d.HasChange("healthchecks") {
		if err = core.Update(resource.Rel("healtchecks"), lbHealthChecks(d)); err != nil {
			return
		}
	}
	if d.HasChange("routingrules") {
		if err = core.Update(resource.Rel("rules"), lbRules(d)); err != nil {
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
