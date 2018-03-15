package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

var firewallSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"description": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"virtualdatacenter": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
	"rules": &schema.Schema{
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"protocol": &schema.Schema{
					Required:     true,
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"ALL", "TCP", "UDP"}, false),
				},
				"fromport": &schema.Schema{
					Required:     true,
					Type:         schema.TypeInt,
					ValidateFunc: validatePort,
				},
				"toport": &schema.Schema{
					Required:     true,
					Type:         schema.TypeInt,
					ValidateFunc: validatePort,
				},
				"targets": &schema.Schema{
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					MinItems: 1,
					Optional: true,
					Type:     schema.TypeList,
				},
				"sources": &schema.Schema{
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					MinItems: 1,
					Optional: true,
					Type:     schema.TypeList,
				},
			},
		},
		MinItems: 1,
		Required: true,
		Type:     schema.TypeList,
	},
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
