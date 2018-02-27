package main

import "github.com/hashicorp/terraform/helper/schema"

func fieldIP() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		ValidateFunc: validateIP,
	}
}

func required(s *schema.Schema) *schema.Schema {
	s.Required = true
	return s
}

func optional(s *schema.Schema) *schema.Schema {
	s.Optional = true
	return s
}

func fieldString() *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeString,
	}
}

func validate(s *schema.Schema, f schema.SchemaValidateFunc) *schema.Schema {
	s.ValidateFunc = f
	return s
}

func fieldTS() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		ValidateFunc: validateTS,
	}
}

func fieldMap(s *schema.Schema) *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeMap,
		Elem: s,
	}
}
