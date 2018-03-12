package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var backupSchema = map[string]*schema.Schema{
	"datacenter": &schema.Schema{
		ForceNew:     true,
		Required:     true,
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	},
	"name": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"code": &schema.Schema{
		Required: true,
		Type:     schema.TypeString,
	},
	"configurations": &schema.Schema{
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"subtype": &schema.Schema{
					Required:     true,
					Type:         schema.TypeString,
					ValidateFunc: validateString([]string{"DEFINED_HOUR", "HOURLY", "DAILY", "MONTHLY", "WEEKLY_PLANNED"}),
				},
				// XXX If date is not properly set in the DTO it generates a GEN-13
				"time": &schema.Schema{
					Default:  "NOT_APPLY",
					Optional: true,
					Type:     schema.TypeString,
				},
				"type": &schema.Schema{
					Required:     true,
					Type:         schema.TypeString,
					ValidateFunc: validateString([]string{"COMPLETE", "SNAPSHOT", "FILESYSTEM"}),
				},
				"days": &schema.Schema{
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validateString([]string{"wednesday", "monday", "tuesday", "thursday", "friday", "saturday", "sunday"}),
					},
					MinItems: 1,
					Optional: true,
					Set:      schema.HashString,
					Type:     schema.TypeSet,
				},
			},
		},
		MinItems: 1,
		Required: true,
		Type:     schema.TypeList,
	},
	"description": &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	},
	"replication": &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	},
}

var backupResource = &schema.Resource{
	Schema: backupSchema,
	Read:   resourceRead(backupDTO, backupRead, "backuppolicy"),
	Update: resourceUpdate(backupDTO, nil, "backuppolicy"),
	Exists: resourceExists("backuppolicy"),
	Delete: resourceDelete,
	Create: resourceCreate(backupDTO, nil, backupRead, backupEndpoint),
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
