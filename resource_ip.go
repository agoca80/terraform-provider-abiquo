package main

import (
	"strings"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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
		ValidateFunc: validation.StringInSlice([]string{"privateip", "externalip", "publicip"}, false),
	},
	"network": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
}

func ipLink(href string) *core.Link {
	var media string
	if private := strings.Contains(href, "/privatenetworks/"); private {
		media = "privateip"
	} else if public := strings.Contains(href, "/publicips/"); public {
		media = "publicip"
	} else {
		media = "externalip"
	}
	return core.NewLinkType(href, media)
}

func ipNew(d *resourceData) core.Resource {
	return &abiquo.IP{IP: d.string("ip")}
}

func ipEndpoint(d *resourceData) *core.Link {
	var media string
	if d.string("type") != "privateip" {
		media = "publicip"
	} else {
		media = "privateip"
	}
	return core.NewLinkType(d.string("network")+"/ips", media)
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
