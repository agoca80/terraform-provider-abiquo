package main

import (
	"strings"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var ipSchema = map[string]*schema.Schema{
	"ip": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateIP,
	},
	"type": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateString([]string{"privateip", "externalip", "publicip"}),
	},
	"network": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
}

func ipLink(href string) *core.Link {
	switch {
	case strings.Contains(href, "/privatenetworks/"):
		return core.NewLinkType(href, "privateip")
	case strings.Contains(href, "/externalnetworks/"):
		return core.NewLinkType(href, "externalip")
	case strings.Contains(href, "/publicips/"):
		return core.NewLinkType(href, "publicip")
	default:
		return nil
	}
}

func ipNew(d *resourceData) core.Resource {
	return &abiquo.IP{IP: d.string("ip")}
}

func ipEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("network")+"/ips", d.string("type"))
}

func ipCreate(rd *schema.ResourceData, meta interface{}) (err error) {
	d := newResourceData(rd, "")
	ip := ipNew(d)
	if err = core.Create(ipEndpoint(d), ip); err == nil {
		d.SetId(ip.URL())
	}
	return
}

// IPResource does not change
func ipRead(rd *schema.ResourceData, meta interface{}) (err error) {
	return core.Read(newResourceData(rd, rd.Get("type").(string)), nil)
}

func ipExists(rd *schema.ResourceData, meta interface{}) (ok bool, err error) {
	err = core.Read(newResourceData(rd, rd.Get("type").(string)), nil)
	return err == nil, err
}
