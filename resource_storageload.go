package main

import (
	"strings"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var storageLoadSchema = map[string]*schema.Schema{
	"load":   Required().Number(),
	"target": Required().Renew().Link(),
}

var storageLoadResource = &schema.Resource{
	Schema: storageLoadSchema,
	Delete: resourceDelete,
	Exists: resourceExists("datastoreloadrule"),
	Create: resourceCreate(storageLoadDTO, nil, storageLoadRead, storageLoadEndpoint),
	Update: resourceUpdate(storageLoadDTO, nil, "datastoreloadrule"),
	Read:   resourceRead(storageLoadDTO, storageLoadRead, "datastoreloadrule"),
}

func storageLoadDTO(d *resourceData) core.Resource {
	storageLoadRule := &abiquo.DatastoreLoadRule{
		StorageLoadPercentage: d.int("load"),
	}

	if h, ok := d.GetOk("target"); ok {
		var media string
		var href = h.(string)
		switch {
		case strings.Contains(href, "datastores"):
			media = "datastore"
		case strings.Contains(href, "datastoretier"):
			media = "datastoretier"
		case strings.Contains(href, "datacenter"):
			media = "datacenter"
		default:
			return nil
		}
		storageLoadRule.Add(core.NewLinkType(href, media).SetRel(media))
	}

	return storageLoadRule
}

func storageLoadEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType("admin/rules/datastoreloadlevel", "datastoreloadrule")
}

func storageLoadRead(d *resourceData, resource core.Resource) (err error) {
	rule := resource.(*abiquo.DatastoreLoadRule)
	d.Set("load", rule.StorageLoadPercentage)
	for _, media := range []string{"datastore", "datastoretier", "datacenter"} {
		if rel := rule.Rel(media); rel != nil {
			d.Set("target", rel.URL())
		}
	}
	return
}
