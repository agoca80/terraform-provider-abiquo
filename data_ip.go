package main

import (
	"fmt"
	"net/url"
	"strings"
	"github.com/hashicorp/terraform/helper/schema"
)

var ipDataSchema = map[string]*schema.Schema{
	"ip":   attribute(required, ip),
	"pool": attribute(required, href),
}

func ipsMedia(pool string) (media string) {
	if strings.Contains(pool, "/privatenetworks/") {
		media = "privateips"
	} else if strings.Contains(pool, "/publicips/") {
		media = "publicips"
	} else {
		media = "externalips"
	}
	return
}

func ipFind(d *resourceData) (err error) {
	address := d.string("ip")
	pool := d.string("pool")
	query := url.Values{"hasIP": {address}}
	ip := linkType(pool, ipsMedia(pool)).Collection(query).First()
	if ip == nil {
		return fmt.Errorf("ip %q not found", address)
	}

	d.SetId(ip.URL())
	return
}
