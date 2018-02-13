package main

import (
	"github.com/abiquo/opal/abiquo"
	"github.com/abiquo/opal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var fwRuleSchema = map[string]*schema.Schema{
	"protocol": Required().ValidateString([]string{"ALL", "TCP", "UDP"}),
	"fromport": Required().Number(),
	"toport":   Required().Number(),
	"targets":  optional(&schema.Schema{Type: schema.TypeList, MinItems: 1, Elem: fieldString()}),
	"sources":  optional(&schema.Schema{Type: schema.TypeList, MinItems: 1, Elem: fieldString()}),
}

var firewallResource = &schema.Resource{
	// Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough},
	Schema: map[string]*schema.Schema{
		"name":              Required().String(),
		"description":       Required().String(),
		"virtualdatacenter": Required().Renew().Link(),
		"rules": required(&schema.Schema{Type: schema.TypeList, MinItems: 1, Elem: &schema.Resource{
			Schema: fwRuleSchema,
		}}),
	},
	Delete: resourceDelete,
	Exists: resourceExists("firewallpolicy"),
	Update: fwUpdate,
	Create: fwCreate,
	Read:   resourceRead(fwNew, fwRead, "firewallpolicy"),
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

func fwEndpoint(d *resourceData) *core.Link {
	vdc := new(abiquo.VirtualDatacenter)
	if core.Read(d.link("virtualdatacenter"), vdc) != nil {
		return nil
	}
	endpoint := vdc.Rel("device")
	if endpoint == nil {
		return nil
	}
	device := new(abiquo.Device)
	if core.Read(endpoint, device) != nil {
		return nil
	}
	return device.Rel("firewalls").SetType("firewallpolicy")
}

func fwCreate(rd *schema.ResourceData, m interface{}) (err error) {
	d := newResourceData(rd, "")
	fw := fwNew(d).(*abiquo.Firewall)
	if err = core.Create(fwEndpoint(d), fw); err == nil {
		d.SetId(fw.URL())
		err = fw.SetRules(fwRules(d))
	}
	return
}

func fwUpdate(rd *schema.ResourceData, m interface{}) (err error) {
	d := newResourceData(rd, "firewallpolicy")
	fw := fwNew(d).(*abiquo.Firewall)
	if err = core.Update(d, fw); err == nil {
		err = fw.SetRules(fwRules(d))
	}
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
