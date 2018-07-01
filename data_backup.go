package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"

	"github.com/hashicorp/terraform/helper/schema"
)

var backupDataSchema = map[string]*schema.Schema{
	"code":     attribute(required, text),
	"location": attribute(required, link("location")),
}

func backupDataRead(d *schema.ResourceData, meta interface{}) (err error) {
	location := d.Get("location").(string)
	endpoint := core.NewLinkType(location, "datacenter")
	resource := endpoint.Walk()
	if resource == nil {
		return fmt.Errorf("location %q does not exist", location)
	}

	code := d.Get("code").(string)
	backups := resource.Rel("backuppolicies").Collection(nil)
	backup := backups.Find(func(r core.Resource) bool {
		return r.(*abiquo.BackupPolicy).Code == code
	})
	if backup == nil {
		return fmt.Errorf("backup %q does not exist in %q", code, location)
	}
	d.SetId(backup.URL())

	return
}
