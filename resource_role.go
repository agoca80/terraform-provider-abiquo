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
	"privileges": attribute(optional, setFn(text, privilegeID)),
}

func roleNew(d *resourceData) core.Resource {
	return &abiquo.Role{
		Name:    d.string("name"),
		Blocked: d.boolean("blocked"),
		DTO:     core.NewDTO(d.link("enterprise")),
	}
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

var role = &description{
	Resource: &schema.Resource{Schema: roleSchema},
	dto:      roleNew,
	endpoint: endpointConst("admin/roles"),
	media:    "role",
	create:   rolePrivileges,
	read:     roleRead,
	update:   rolePrivileges,
}
