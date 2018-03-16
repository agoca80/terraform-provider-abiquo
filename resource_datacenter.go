package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var datacenterSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"location": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"vf": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"vsm": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"am": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"nc": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"ssm": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"bpm": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"cpp": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"dhcp": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"dhcpv6": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"ra": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
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
	rss := make([]abiquo.RemoteService, len(rssMap))
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
	for _, rs := range datacenter.RemoteServices.Collection {
		d.Set(rssMap[rs.Type], rs.URI)
	}
	return
}
