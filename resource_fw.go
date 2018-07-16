package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var fwRulesResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"protocol": attribute(required, protocol),
		"fromport": attribute(required, port),
		"toport":   attribute(required, port),
		"targets":  attribute(optional, list(text), min(1)),
		"sources":  attribute(optional, list(text), min(1)),
	},
}

var firewallSchema = map[string]*schema.Schema{
	"virtualdatacenter": endpoint("virtualdatacenter"),
	"name":              attribute(required, text),
	"description":       attribute(required, text),
	"rules":             attribute(required, list(fwRulesResource), min(1)),
}

func fwNew(d *resourceData) core.Resource {
	return &abiquo.Firewall{
		Name:        d.string("name"),
		Description: d.string("description"),
		DTO: core.NewDTO(
			d.link("virtualdatacenter"),
		),
	}
}

func fwRules(d *resourceData) *abiquo.FirewallRules {
	slice := d.slice("rules")
	rules := make([]abiquo.FirewallRule, len(slice))
	for i, r := range slice {
		mapDecoder(r, &rules[i])
	}
	return &abiquo.FirewallRules{
		Collection: rules,
	}
}

func fwEndpoint(d *resourceData) (link *core.Link) {
	if device := vdcDevice(d.link("virtualdatacenter")); device != nil {
		link = device.Rel("firewalls").SetType("firewallpolicy")
	}
	return
}

func fwCreate(d *resourceData, resource core.Resource) (err error) {
	fw := resource.(*abiquo.Firewall)
	if rules := fwRules(d); len(rules.Collection) > 0 {
		err = fw.SetRules(fwRules(d))
	}
	return
}

func fwUpdate(d *resourceData, resource core.Resource) (err error) {
	fw := resource.(*abiquo.Firewall)
	err = fw.SetRules(fwRules(d))
	return
}

func fwRead(d *resourceData, resource core.Resource) (err error) {
	// Read the firewall
	fw := resource.(*abiquo.Firewall)
	d.Set("name", fw.Name)
	d.Set("description", fw.Description)

	// Read the firewall rules
	rules := new(abiquo.FirewallRules)
	if err = core.Read(fw.Rel("rules"), rules); err != nil {
		return
	}

	value := make([]interface{}, len(rules.Collection))
	for i, r := range rules.Collection {
		value[i] = map[string]interface{}{
			"fromport": r.FromPort,
			"toport":   r.ToPort,
			"protocol": r.Protocol,
			"sources":  r.Sources,
			"targets":  r.Targets,
		}
	}
	d.Set("rules", value)
	return
}
