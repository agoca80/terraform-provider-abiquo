package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var roleSchema = map[string]*schema.Schema{
	"blocked":    attribute(boolean, optional),
	"name":       attribute(required, text),
	"enterprise": attribute(optional, link("enterprise")),
	"privileges": &schema.Schema{
		Elem:     attribute(text),
		Required: true,
		Set:      privilegeID,
		Type:     schema.TypeSet,
	},
}

func roleNew(d *resourceData) core.Resource {
	return &abiquo.Role{
		Name:    d.string("name"),
		Blocked: d.boolean("blocked"),
		DTO:     core.NewDTO(d.link("enterprise")),
	}
}

func roleEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType("admin/roles", "role")
}

func roleRead(d *resourceData, resource core.Resource) (err error) {
	role := resource.(*abiquo.Role)
	privileges := schema.NewSet(privilegeID, nil)
	collection := role.Rel("privileges").Collection(nil)
	for _, p := range collection.List() {
		privileges.Add(p.(*abiquo.Privilege).Name)
	}
	d.Set("name", role.Name)
	d.Set("privileges", privileges)
	if _, ok := d.GetOk("blocked"); ok {
		d.Set("blocked", role.Blocked)
	}
	return
}

func rolePrivileges(d *resourceData, resource core.Resource) (err error) {
	if !d.HasChange("privileges") {
		return
	}

	role := resource.(*abiquo.Role)
	for _, name := range d.set("privileges").List() {
		privilege := privilegeGet(name.(string))
		if privilege == nil {
			return fmt.Errorf("roleCreate: privilege %v does not exist", name)
		}
		role.AddPrivilege(privilege)
	}
	return core.Update(role, role)
}

var resourceRole = &schema.Resource{
	Schema: roleSchema,
	Delete: resourceDelete,
	Read:   resourceRead(roleNew, roleRead, "role"),
	Create: resourceCreate(roleNew, rolePrivileges, roleRead, roleEndpoint),
	Exists: resourceExists("role"),
	Update: resourceUpdate(roleNew, rolePrivileges, "role"),
}
