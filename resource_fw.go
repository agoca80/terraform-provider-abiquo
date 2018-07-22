package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var firewallSchema = map[string]*schema.Schema{
	"device":            endpoint("device"),
	"virtualdatacenter": attribute(required, link("virtualdatacenter"), forceNew),
	"name":              attribute(required, text),
	"description":       attribute(required, text),
	"rules": attribute(required, min(1), list(&schema.Resource{
		Schema: map[string]*schema.Schema{
			"protocol": attribute(required, protocol),
			"fromport": attribute(required, port),
			"toport":   attribute(required, port),
			"targets":  attribute(optional, list(text), min(1)),
			"sources":  attribute(optional, list(text), min(1)),
		},
	})),
}

func firewallpolicyNew(d *resourceData) core.Resource {
	return &abiquo.Firewall{
		Name:        d.string("name"),
		Description: d.string("description"),
		DTO: core.NewDTO(
			d.link("virtualdatacenter"),
		),
	}
}

func firewallpolicyRules(d *resourceData) *abiquo.FirewallRules {
	slice := d.slice("rules")
	rules := make([]abiquo.FirewallRule, len(slice))
	for i, r := range slice {
		mapDecoder(r, &rules[i])
	}
	return &abiquo.FirewallRules{
		Collection: rules,
	}
}

func firewallpolicyUpdateRules(d *resourceData, resource core.Resource) (err error) {
	firewallpolicy := resource.(*abiquo.Firewall)
	if d.HasChange("rules") {
		err = core.Update(firewallpolicy.Rel("rules"), firewallpolicyRules(d))
	}
	return
}

func firewallpolicyRead(d *resourceData, resource core.Resource) (err error) {
	// Read the firewall
	firewallpolicy := resource.(*abiquo.Firewall)
	d.Set("name", firewallpolicy.Name)
	d.Set("description", firewallpolicy.Description)

	// Read the firewall rules
	rules := new(abiquo.FirewallRules)
	if err = core.Read(firewallpolicy.Rel("rules"), rules); err != nil {
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

var firewallpolicy = &description{
	media:    "firewallpolicy",
	dto:      firewallpolicyNew,
	read:     firewallpolicyRead,
	create:   firewallpolicyUpdateRules,
	update:   firewallpolicyUpdateRules,
	endpoint: endpointPath("device", "/firewalls"),
	Resource: &schema.Resource{Schema: firewallSchema},
}
