package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var datacenterSchema = map[string]*schema.Schema{
	"name":     attribute(required, text),
	"location": attribute(required, text),
	"vf":       attribute(required, text),
	"vsm":      attribute(required, text),
	"am":       attribute(required, text),
	"nc":       attribute(required, text),
	"ssm":      attribute(required, text),
	"bpm":      attribute(required, text),
	"cpp":      attribute(required, text),
	"dhcp":     attribute(required, text),
	"dhcpv6":   attribute(required, text),
	"ra":       attribute(required, text),
	"tiers":    attribute(computed, text),
}

var rssMap = map[string]string{
	"VIRTUAL_FACTORY":        "vf",
	"VIRTUAL_SYSTEM_MONITOR": "vsm",
	"APPLIANCE_MANAGER":      "am",
	"NODE_COLLECTOR":         "nc",
	"STORAGE_SYSTEM_MONITOR": "ssm",
	"BPM_SERVICE":            "bpm",
	"CLOUD_PROVIDER_PROXY":   "cpp",
	"DHCP_SERVICE":           "dhcp",
	"DHCPv6":                 "dhcpv6",
	"REMOTE_ACCESS":          "ra",
}

func datacenterNew(d *resourceData) core.Resource {
	rss := []abiquo.RemoteService{}
	for k, v := range rssMap {
		rss = append(rss, abiquo.RemoteService{Type: k, URI: d.string(v)})
	}
	datacenter := &abiquo.Datacenter{
		Name:     d.string("name"),
		Location: d.string("location"),
	}
	datacenter.RemoteServices.Collection = rss
	return datacenter
}

func datacenterEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType("admin/datacenters", "datacenter")
}

func datacenterRead(d *resourceData, resource core.Resource) (err error) {
	datacenter := resource.(*abiquo.Datacenter)
	d.Set("name", datacenter.Name)
	d.Set("location", datacenter.Location)
	d.Set("tiers", datacenter.Rel("tiers").Href)
	for _, rs := range datacenter.RemoteServices.Collection {
		d.Set(rssMap[rs.Type], rs.URI)
	}
	return
}
