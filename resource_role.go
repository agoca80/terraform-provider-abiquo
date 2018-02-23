package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var roleResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"blocked":    Optional().Bool(),
		"name":       Required().String(),
		"enterprise": Optional().Link(),
		"privileges": optional(fieldSet(privilegeID, 1, schema.TypeString)),
	},
	Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough},
	Delete:   resourceDelete,
	Read:     resourceRead(roleNew, roleRead, "role"),
	Create:   roleCreate,
	Exists:   resourceExists("role"),
	Update:   roleUpdate,
}

func roleNew(d *resourceData) core.Resource {
	return &abiquo.Role{
		Name:    d.string("name"),
		Blocked: d.bool("blocked"),
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

func roleCreate(rd *schema.ResourceData, m interface{}) (err error) {
	d := newResourceData(rd, "role")
	role := roleNew(d).(*abiquo.Role)
	if err = core.Create(roleEndpoint(nil), role); err != nil {
		return
	}
	d.SetId(role.URL())
	return rolePrivilegesUpdate(role, d)
}

func roleUpdate(rd *schema.ResourceData, m interface{}) (err error) {
	d := newResourceData(rd, "role")
	role := roleNew(d).(*abiquo.Role)
	if err = core.Update(d, role); err == nil {
		err = rolePrivilegesUpdate(role, d)
	}
	return
}

func rolePrivilegesUpdate(r *abiquo.Role, d *resourceData) error {
	for _, name := range d.set("privileges").List() {
		privilege := privilegeGet(name.(string))
		if privilege == nil {
			return fmt.Errorf("roleCreate: privilege %v does not exist", name)
		}
		r.AddPrivilege(privilege)
	}
	return core.Update(r, r)
}
