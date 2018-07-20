package main

import (
	"fmt"
	"sync"

	"github.com/abiquo/ojal/abiquo"
	"github.com/hashicorp/terraform/helper/schema"
)

var privileges = struct {
	sync.Once
	privilege map[string]*abiquo.Privilege
}{}

var privilegeDataSchema = map[string]*schema.Schema{
	"name": attribute(required, text),
}

func privilegeRead(d *schema.ResourceData, meta interface{}) (err error) {
	privilege := privilegeGet(d.Get("name").(string))
	if privilege == nil {
		return fmt.Errorf("Privilege %v does not exist", d.Get("name"))
	}
	d.SetId(privilege.URL())
	d.Set("name", privilege.Name)
	return
}

func privilegeGet(name string) *abiquo.Privilege {
	privileges.Do(func() {
		privileges.privilege = make(map[string]*abiquo.Privilege)
		for _, p := range abiquo.Privileges(nil).List() {
			privilege := p.(*abiquo.Privilege)
			privileges.privilege[privilege.Name] = privilege
		}
	})
	return privileges.privilege[name]
}

func privilegeID(name interface{}) (id int) {
	if p := privilegeGet(name.(string)); p != nil {
		id = p.ID
	}
	return
}

var dataPrivilege = &schema.Resource{
	Schema: privilegeDataSchema,
	Read:   privilegeRead,
}
