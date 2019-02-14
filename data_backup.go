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

func backupFind(d *resourceData) (err error) {
	location := d.string("location")
	endpoint := linkType(location, "datacenter")
	resource, err := endpoint.Walk()
	if err != nil {
		return
	}

	code := d.string("code")
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
