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

func datacenterNew(d *resourceData) core.Resource {
	rss := make([]abiquo.RemoteService, 10)
	rss[0] = abiquo.RemoteService{Type: "VIRTUAL_FACTORY", URI: d.string("vf")}
	rss[1] = abiquo.RemoteService{Type: "VIRTUAL_SYSTEM_MONITOR", URI: d.string("vsm")}
	rss[2] = abiquo.RemoteService{Type: "APPLIANCE_MANAGER", URI: d.string("am")}
	rss[3] = abiquo.RemoteService{Type: "NODE_COLLECTOR", URI: d.string("nc")}
	rss[4] = abiquo.RemoteService{Type: "STORAGE_SYSTEM_MONITOR", URI: d.string("ssm")}
	rss[5] = abiquo.RemoteService{Type: "BPM_SERVICE", URI: d.string("bpm")}
	rss[6] = abiquo.RemoteService{Type: "CLOUD_PROVIDER_PROXY", URI: d.string("cpp")}
	rss[7] = abiquo.RemoteService{Type: "DHCP_SERVICE", URI: d.string("dhcp")}
	rss[8] = abiquo.RemoteService{Type: "DHCPv6", URI: d.string("dhcpv6")}
	rss[9] = abiquo.RemoteService{Type: "REMOTE_ACCESS", URI: d.string("ra")}
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
		switch rs.Type {
		case "VIRTUAL_FACTORY":
			d.Set("vf", rs.URI)
		case "VIRTUAL_SYSTEM_MONITOR":
			d.Set("vsm", rs.URI)
		case "APPLIANCE_MANAGER":
			d.Set("am", rs.URI)
		case "NODE_COLLECTOR":
			d.Set("nc", rs.URI)
		case "STORAGE_SYSTEM_MONITOR":
			d.Set("ssm", rs.URI)
		case "BPM_SERVICE":
			d.Set("bpm", rs.URI)
		case "CLOUD_PROVIDER_PROXY":
			d.Set("cpp", rs.URI)
		case "DHCP_SERVICE":
			d.Set("dhcp", rs.URI)
		case "DHCPv6":
			d.Set("dhcpv6", rs.URI)
		case "REMOTE_ACCESS":
			d.Set("ra", rs.URI)
		}
	}
	return
}
