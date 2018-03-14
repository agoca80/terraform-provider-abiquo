package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var roleSchema = map[string]*schema.Schema{
	"blocked": &schema.Schema{
		Optional: true,
		Type:     schema.TypeBool,
	},
	"name": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"enterprise": &schema.Schema{
		Optional:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
	"privileges": &schema.Schema{
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Required: true,
		Set:      privilegeID,
		Type:     schema.TypeSet,
	},
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
