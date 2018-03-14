package main

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/abiquo/ojal/abiquo"
	"github.com/hashicorp/terraform/helper/schema"
)

var machineDataSchema = map[string]*schema.Schema{
	// Discover parameters
	"hypervisor": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"datacenter": &schema.Schema{
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
	"ip": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"port": &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	},
	"managerip": &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	},
	"manageruser": &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	},
	"managerpass": &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	},
	"definition": &schema.Schema{
		Computed: true,
		Type:     schema.TypeString,
	},
}

var machineDataSource = &schema.Resource{
	Schema: machineDataSchema,
	Read:   machineDataRead,
}

func machineDataRead(rd *schema.ResourceData, _ interface{}) (err error) {
	d := newResourceData(rd, "")

	var query url.Values
	hypervisor := d.string("hypervisor")
	switch hypervisor {
	case "KVM":
		query = url.Values{
			"ip": {d.string("ip")},
		}
	case "VMX_04":
		query = url.Values{
			"ip":              {d.string("ip")},
			"managerip":       {d.string("managerip")},
			"manageruser":     {d.string("manageruser")},
			"managerpassword": {d.string("managerpass")},
		}
	default:
		return fmt.Errorf("unknown hypervisor type: %q", hypervisor)
	}

	resource := d.link("datacenter").SetType("datacenter").Walk()
	if resource == nil {
		return fmt.Errorf("datacenter not found: %q", d.string("datacenter"))
	}
	datacenter := resource.(*abiquo.Datacenter)

	query["hypervisor"] = []string{hypervisor}
	if port := d.string("port"); port != "" {
		query["port"] = []string{port}
	}
	resource = datacenter.Rel("discover").Collection(query).First()
	if resource == nil {
		return fmt.Errorf("machine not found: %v", query)
	}

	bytes, err := json.Marshal(resource)
	if err != nil {
		return
	}
	d.SetId(d.string("ip"))
	d.Set("definition", string(bytes))
	return
}
