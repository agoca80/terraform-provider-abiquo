package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func optional(s *schema.Schema)  { s.Optional = true }
func computed(s *schema.Schema)  { s.Computed = true }
func required(s *schema.Schema)  { s.Required = true }
func forceNew(s *schema.Schema)  { s.ForceNew = true }
func sensitive(s *schema.Schema) { s.Sensitive = true }

func text(s *schema.Schema)    { s.Type = schema.TypeString }
func integer(s *schema.Schema) { s.Type = schema.TypeInt }
func boolean(s *schema.Schema) { s.Type = schema.TypeBool }

func hash(elem interface{}) field {
	return func(s *schema.Schema) {
		s.Elem = elem
		s.Type = schema.TypeMap
	}
}

func set(elem interface{}, set schema.SchemaSetFunc) field {
	return func(s *schema.Schema) {
		s.Elem = elem
		s.Set = set
		s.Type = schema.TypeSet
	}
}

func setLink(media string) field { return set(attribute(link(media)), schema.HashString) }

func list(elem interface{}) field {
	return func(s *schema.Schema) {
		s.Elem = elem
		s.Type = schema.TypeList
	}
}

func min(m int) field {
	return func(s *schema.Schema) {
		s.MinItems = m
	}
}

func port(s *schema.Schema) {
	integer(s)
	s.ValidateFunc = validatePort
}

func protocol(s *schema.Schema) {
	text(s)
	s.ValidateFunc = validation.StringInSlice([]string{"TCP", "HTTP", "HTTPS"}, false)
}

func price(s *schema.Schema) {
	s.Type = schema.TypeFloat
	s.Optional = true
	s.ValidateFunc = validatePrice
}

func natural(s *schema.Schema) {
	atLeast(0)(s)
}

func atLeast(m int) field {
	return func(s *schema.Schema) {
		integer(s)
		s.ValidateFunc = validation.IntAtLeast(m)
	}
}

func ip(s *schema.Schema) {
	text(s)
	s.ValidateFunc = validateIP
}

func timestamp(s *schema.Schema) {
	text(s)
	s.ValidateFunc = validateTS
}

func href(s *schema.Schema) {
	text(s)
	s.ValidateFunc = validateHref
}

func link(media string) field {
	return func(s *schema.Schema) {
		text(s)
		s.ValidateFunc = validateMedia[media]
	}
}

func label(strs []string) field {
	return func(s *schema.Schema) {
		text(s)
		s.ValidateFunc = validation.StringInSlice(strs, false)
	}
}

func conflicts(strs []string) field {
	return func(s *schema.Schema) {
		s.ConflictsWith = strs
	}
}

type field func(*schema.Schema)

func attribute(fields ...field) (media *schema.Schema) {
	media = &schema.Schema{}
	for _, field := range fields {
		field(media)
	}
	return
}

func resourceSet(v interface{}) int {
	resource := v.(map[string]interface{})
	return schema.HashString(resource["href"].(string))
}

func variable(name string) field {
	return func(s *schema.Schema) {
		s.DefaultFunc = schema.EnvDefaultFunc(name, "")
	}
}

func prices(s *schema.Schema) {
	s.Set = resourceSet
	s.Type = schema.TypeSet
	s.Elem = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"href":  attribute(required, href),
			"price": attribute(price),
		},
	}
}
