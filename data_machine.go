package main

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
)

var machineDataSchema = map[string]*schema.Schema{
	"hypervisor":  attribute(required, label(machineType)),
	"datacenter":  attribute(required, link("datacenter")),
	"ip":          attribute(required, ip),
	"port":        attribute(optional, text),
	"managerip":   attribute(optional, text),
	"manageruser": attribute(optional, text),
	"managerpass": attribute(optional, text),
	"definition":  attribute(computed, text),
}

func machineFind(d *resourceData) (err error) {
	var query url.Values
	switch d.string("hypervisor") {
	case "KVM":
		query = url.Values{
			"ip":         {d.string("ip")},
			"hypervisor": {"KVM"},
		}
	case "VMX_04":
		query = url.Values{
			"ip":              {d.string("ip")},
			"managerip":       {d.string("managerip")},
			"manageruser":     {d.string("manageruser")},
			"managerpassword": {d.string("managerpass")},
			"hypervisor":      {"VMX_04"},
		}
	}

	datacenter := d.link("datacenter").SetType("datacenter").Walk()
	if datacenter == nil {
		return fmt.Errorf("datacenter not found: %q", d.string("datacenter"))
	}

	if port := d.string("port"); port != "" {
		query["port"] = []string{port}
	}
	resource := datacenter.Rel("discover").Collection(query).First()
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
