package main

import (
	"strings"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var ipSchema = map[string]*schema.Schema{
	"ip":      attribute(required, ip, forceNew),
	"type":    attribute(computed, text),
	"network": attribute(required, href, forceNew),
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
	return linkType(href, media)
}

func ipCreate(rd *schema.ResourceData, meta interface{}) (err error) {
	href := rd.Get("network").(string) + "/ips"

	var media string
	if private := strings.Contains(href, "/privatenetworks/"); private {
		media = "privateip"
	} else {
		media = "publicip"
	}

	ip := &abiquo.IP{
		IP:        rd.Get("ip").(string),
		Available: true,
	}

	err = linkType(href, media).Create(ip)
	if err != nil {
		return
	}

	rd.SetId(ip.URL())
	rd.Set("type", ip.Media())
	return
}

// IPResource does not change
func ipRead(rd *schema.ResourceData, meta interface{}) (err error) {
	href := rd.Id()
	media := rd.Get("type").(string)
	endpoint := linkType(href, media)
	err = endpoint.Read(nil)
	return
}

func ipExists(rd *schema.ResourceData, meta interface{}) (ok bool, err error) {
	href := rd.Id()
	media := rd.Get("type").(string)
	endpoint := linkType(href, media)
	err = endpoint.Read(nil)
	return err == nil, err
}

var ipAddress = &description{
	name: "ip",
	Resource: &schema.Resource{
		Schema: ipSchema,
		Exists: ipExists,
		Create: ipCreate,
		Read:   ipRead,
	},
}
