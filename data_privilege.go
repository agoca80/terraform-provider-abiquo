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

func privilegeFind(d *resourceData) (err error) {
	privilege := privilegeGet(d.string("name"))
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
