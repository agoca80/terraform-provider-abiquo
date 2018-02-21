package main

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/abiquo/opal/abiquo"
	"github.com/hashicorp/terraform/helper/schema"
)

var machineDataSchema = map[string]*schema.Schema{
	// Discover parameters
	"hypervisor":  Required().String(),
	"datacenter":  Required().ValidateURL(),
	"ip":          Required().String(),
	"managerip":   Optional().String(),
	"manageruser": Optional().String(),
	"managerpass": Optional().String(),
	"definition":  Computed().String(),
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
	resource = datacenter.Discover(query).First()
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
