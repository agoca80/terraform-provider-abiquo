package main

import (
	"encoding/json"
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var machineType = []string{"VMX_04", "KVM"}

var machineInterface = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name": attribute(required, text),
		"nst":  attribute(required, href),
	},
}

func interfaceSet(v interface{}) int {
	return schema.HashString(v.(map[string]interface{})["name"])
}

var machineDatastore = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"uuid":   attribute(required, text),
		"dstier": attribute(required, href),
	},
}

func datastoreSet(v interface{}) int {
	return schema.HashString(v.(map[string]interface{})["uuid"])
}

var machineSchema = map[string]*schema.Schema{
	"rack":        endpoint("rack"),
	"definition":  attribute(required, text),
	"datastore":   attribute(required, setFn(machineDatastore, datastoreSet), min(1)),
	"interface":   attribute(required, setFn(machineInterface, interfaceSet), min(1)),
	"managerip":   attribute(optional, ip),
	"manageruser": attribute(optional, text),
	"managerpass": attribute(optional, text, sensitive),
	"type":        attribute(computed, text),
}

func machineCreate(rd *schema.ResourceData, _ interface{}) (err error) {
	d := newDataType(rd, "machine")
	definition := d.string("definition")
	machine := new(abiquo.Machine)
	if err = json.Unmarshal([]byte(definition), machine); err != nil {
		return fmt.Errorf("definition is not a valid machine: %q", definition)
	}

	if machine.Type == "VMX_04" {
		machine.ManagerIP = d.string("managerip")
		machine.ManagerUser = d.string("manageruser")
		machine.ManagerPass = d.string("managerpass")
	}

	interfaces := make(map[string]interface{})
	for _, i := range d.set("interface").List() {
		iface := i.(map[string]interface{})
		interfaces[iface["name"].(string)] = iface["nst"]
	}

	for _, i := range machine.Interfaces.Collection {
		if href, ok := interfaces[i.Name]; ok {
			nst := core.NewLinkType(href.(string), "networkservicetype")
			i.Add(nst.SetRel("networkservicetype"))
		}
	}

	datastores := make(map[string]interface{})
	for _, d := range d.set("datastore").List() {
		datastore := d.(map[string]interface{})
		datastores[datastore["uuid"].(string)] = datastore["dstier"]
	}

	for _, d := range machine.Datastores.Collection {
		if href, ok := datastores[d.UUID]; ok {
			dstier := core.NewLinkType(href.(string), "datastoretier")
			d.Add(dstier.SetRel("datastoretier"))
			d.Enabled = true
		}
	}

	endpoint := core.NewLinkType(d.string("rack")+"/machines", "machine")
	if err = core.Create(endpoint, machine); err == nil {
		d.SetId(machine.URL())
		d.Set("type", machine.Type)
	}
	return
}

var machine = &description{
	media: "machine",
	Resource: &schema.Resource{
		Schema: machineSchema,
		Create: machineCreate,
		Update: schema.Noop,
		Read:   schema.Noop,
	},
}
