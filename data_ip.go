package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/abiquo/ojal/core"

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

func ipDataRead(d *schema.ResourceData, meta interface{}) (err error) {
	address := d.Get("ip").(string)
	pool := d.Get("pool").(string)
	query := url.Values{"hasIP": {address}}
	ip := core.NewLinkType(pool, ipsMedia(pool)).Collection(query).First()
	if ip == nil {
		return fmt.Errorf("ip %q not found", address)
	}

	d.SetId(ip.URL())
	return
}

var dataIp = &schema.Resource{
	Schema: ipDataSchema,
	Read:   ipDataRead,
}
