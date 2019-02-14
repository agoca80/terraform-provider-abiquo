package main

import (
	"strings"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var storageLoadSchema = map[string]*schema.Schema{
	"load":   attribute(required, integer),
	"target": attribute(required, href, forceNew),
}

func storageLoadDTO(d *resourceData) core.Resource {
	storageLoadRule := &abiquo.DatastoreLoadRule{
		StorageLoadPercentage: d.integer("load"),
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
		storageLoadRule.Add(linkType(href, media).SetRel(media))
	}

	return storageLoadRule
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

var datastoreloadrule = &description{
	Resource: &schema.Resource{Schema: storageLoadSchema},
	dto:      storageLoadDTO,
	endpoint: endpointConst("admin/rules/datastoreloadlevel"),
	media:    "datastoreloadrule",
	read:     storageLoadRead,
}
