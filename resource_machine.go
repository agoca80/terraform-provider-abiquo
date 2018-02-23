package main

import (
	"encoding/json"
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var machineSchema = map[string]*schema.Schema{
	"definition": Required().String(),
	"datastores": &schema.Schema{
		Elem:     &schema.Schema{Type: schema.TypeString},
		Required: true,
		Type:     schema.TypeMap,
	},
	"interfaces": &schema.Schema{
		Elem:     &schema.Schema{Type: schema.TypeString},
		Required: true,
		Type:     schema.TypeMap,
	},
	"managerip":   Optional().String(),
	"manageruser": Optional().String(),
	"managerpass": Optional().String(),
	"rack":        Required().ValidateURL(),
}

var machineResource = &schema.Resource{
	Schema: machineSchema,
	Delete: resourceDelete,
	Exists: resourceExists("machine"),
	Create: machineCreate,
	Update: machineUpdate,
	Read:   machineRead,
}

func machineCreate(rd *schema.ResourceData, _ interface{}) (err error) {
	d := newResourceData(rd, "machine")
	definition := d.string("definition")
	machine := new(abiquo.Machine)
	if err = json.Unmarshal([]byte(definition), machine); err != nil {
		return fmt.Errorf("definition is not a valid machine: %q", definition)
	}

	switch machine.Type {
	case "VMX_04":
		machine.ManagerIP = d.string("managerip")
		machine.ManagerUser = d.string("manageruser")
		machine.ManagerPass = d.string("managerpass")
	}

	// Enable interfaces
	ifaces := d.dict("interfaces")
	for _, iface := range machine.Interfaces.Collection {
		if href, ok := ifaces[iface.MAC]; ok {
			nst := core.NewLinkType(href.(string), "networkservicetype")
			iface.Add(nst.SetRel("networkservicetype"))
		}
	}

	// Enable datastores
	dstores := d.dict("datastores")
	for _, dstore := range machine.Datastores.Collection {
		var dstier *core.Link
		href, ok := dstores[dstore.DatastoreUUID]
		if !ok {
			continue
		}
		if href.(string) != "" {
			dstier = core.NewLinkType(href.(string), "datastoretier").SetRel("datastoretier")
		}
		dstore.Enabled = true
		dstore.Add(dstier)
	}

	endpoint := core.NewLinkType(d.string("rack")+"/machines", "machine")
	if err = core.Create(endpoint, machine); err != nil {
		return
	}
	d.SetId(machine.URL())
	return
}

func machineUpdate(rd *schema.ResourceData, _ interface{}) (err error) {
	return
}

func machineRead(rd *schema.ResourceData, _ interface{}) (err error) {
	return
}
