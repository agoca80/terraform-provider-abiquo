package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var backupType = []string{"COMPLETE", "SNAPSHOT", "FILESYSTEM"}
var backupSubtype = []string{"DEFINED_HOUR", "HOURLY", "DAILY", "MONTHLY", "WEEKLY_PLANNED"}
var weekDays = []string{"wednesday", "monday", "tuesday", "thursday", "friday", "saturday", "sunday"}

// XXX If date is not properly set in the DTO it generates a GEN-13
var backupConfiguration = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"subtype": attribute(required, label(backupSubtype)),
		"time":    attribute(optional, text),
		"type":    attribute(required, label(backupType)),
		"days":    attribute(optional, set(attribute(label(weekDays)), schema.HashString), min(1)),
	},
}

var backupSchema = map[string]*schema.Schema{
	"datacenter":     attribute(required, forceNew, link("datacenter")),
	"name":           attribute(required, text),
	"code":           attribute(required, text),
	"configurations": attribute(required, list(backupConfiguration), min(1)),
	"description":    attribute(optional, text),
	"replication":    attribute(optional, text),
}

func backupDTO(d *resourceData) core.Resource {
	confs := []abiquo.BackupConfiguration{}
	for _, value := range d.slice("configurations") {
		conf := value.(map[string]interface{})
		confs = append(confs, abiquo.BackupConfiguration{
			Subtype: conf["subtype"].(string),
			Time:    conf["time"].(string),
			Type:    conf["type"].(string),
		})
	}
	return &abiquo.BackupPolicy{
		Name:           d.string("name"),
		Code:           d.string("code"),
		Configurations: confs,
	}
}

func backupEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType(d.string("datacenter")+"/backuppolicies", "backuppolicy")
}

func backupRead(d *resourceData, resource core.Resource) (err error) {
	backup := resource.(*abiquo.BackupPolicy)
	confs := make([]interface{}, len(backup.Configurations))
	for i, c := range backup.Configurations {
		confs[i] = map[string]interface{}{
			"subtype": c.Subtype,
			"time":    c.Time,
			"type":    c.Type,
		}
	}
	d.Set("confs", confs)
	d.Set("code", backup.Code)
	d.Set("name", backup.Name)
	return
}
