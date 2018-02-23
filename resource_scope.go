package main

import (
	"path"
	"strconv"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var scopeResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name":        Required().String(),
		"parent":      Optional().Link(),
		"datacenters": Optional().Links(),
		"enterprises": Optional().Links(),
	},
	Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough},
	Delete:   resourceDelete,
	Read:     resourceRead(scopeNew, scopeRead, "scope"),
	Create:   resourceCreate(scopeNew, nil, scopeRead, scopeEndpoint),
	Exists:   resourceExists("scope"),
	Update:   resourceUpdate(scopeNew, "scope"),
}

func scopeNew(d *resourceData) core.Resource {
	link := core.NewLinkType("admin/scopes/undefined", "scope").SetRel("scope").SetTitle(d.string("name"))
	entities := []abiquo.ScopeEntity{}
	// Add datacenters to scope Entities
	for i, href := range d.slice("datacenters") {
		id := path.Base(href.(string))
		idResource, err := strconv.Atoi(id)
		if err != nil {
			panic("scopeNew: Unexpected datacenter href format: " + href.(string))
		}

		entities = append(entities, abiquo.ScopeEntity{
			ID:         i + 1, // scope entities IDs start at 1
			IDResource: idResource,
			EntityType: "DATACENTER",
			DTO:        core.NewDTO(link),
		})
	}

	// Add enterprises to scope entities
	for i, v := range d.slice("enterprises") {
		href := v.(string)
		id := path.Base(href)
		idResource, err := strconv.Atoi(id)
		if err != nil {
			panic("scopeNew: Unexpected enterprise href format: " + href)
		}

		entities = append(entities, abiquo.ScopeEntity{
			ID:         i + 1,
			IDResource: idResource,
			EntityType: "ENTERPRISE",
			DTO:        core.NewDTO(link),
		})
	}

	return &abiquo.Scope{
		Name:     d.string("name"),
		Entities: entities,
		DTO: core.NewDTO(
			d.link("parent").SetType("scope").SetRel("scopeParent"),
		)}
}

func scopeEndpoint(d *resourceData) *core.Link {
	return core.NewLinkType("admin/scopes", "scope")
}

func scopeRead(d *resourceData, resource core.Resource) (err error) {
	scope := resource.(*abiquo.Scope)
	datacenters := []string{}
	enterprises := []string{}
	for _, entity := range scope.Entities {
		switch entity.EntityType {
		case "DATACENTER":
			path := "admin/datacenters/" + strconv.Itoa(entity.IDResource)
			href := core.Resolve(path, nil)
			datacenters = append(datacenters, href)
		case "ENTERPRISE":
			path := "admin/enterprises/" + strconv.Itoa(entity.IDResource)
			href := core.Resolve(path, nil)
			enterprises = append(enterprises, href)
		default:
			panic("Illegal scopeEntity type")
		}
	}
	d.Set("name", scope.Name)
	d.Set("datacenters", datacenters)
	d.Set("enterprises", enterprises)
	return
}
