package main

import "github.com/hashicorp/terraform/helper/schema"

func fieldSet(set schema.SchemaSetFunc, min int, kind schema.ValueType) (s *schema.Schema) {
	return &schema.Schema{
		Set:      set,
		Type:     schema.TypeSet,
		MinItems: min,
		Elem: &schema.Schema{
			Type: kind,
		},
	}
}

func fieldStrings() (s *schema.Schema) {
	return &schema.Schema{
		MinItems: 1,
		Type:     schema.TypeList,
		Elem:     fieldString(),
	}
}

func fieldLinks() (s *schema.Schema) {
	return &schema.Schema{
		MinItems: 1,
		Type:     schema.TypeList,
		Elem:     fieldLink(),
	}
}

func fieldIP() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		ValidateFunc: validateIP,
	}
}

func fieldInt() *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeInt,
	}
}

func field(valueType schema.ValueType) *schema.Schema {
	return &schema.Schema{
		Type: valueType,
	}
}

func computed(s *schema.Schema) *schema.Schema {
	s.Computed = true
	return s
}

func required(s *schema.Schema) *schema.Schema {
	s.Required = true
	return s
}

func optional(s *schema.Schema) *schema.Schema {
	s.Optional = true
	return s
}

func fieldBool() *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeBool,
	}
}

func fieldString() *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeString,
	}
}

func fieldLink() (s *schema.Schema) {
	return &schema.Schema{
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	}
}

func fieldCtrlType() (s *schema.Schema) {
	return &schema.Schema{
		Type:         schema.TypeString,
		ValidateFunc: validateController,
	}
}

func validate(s *schema.Schema, f schema.SchemaValidateFunc) *schema.Schema {
	s.ValidateFunc = f
	return s
}

func New(s *schema.Schema) *schema.Schema {
	s.ForceNew = true
	return s
}

func fieldTS() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		ValidateFunc: validateTS,
	}
}

func conflicts(s *schema.Schema, conflict []string) *schema.Schema {
	s.ConflictsWith = conflict
	return optional(s)
}

func fieldMap(s *schema.Schema) *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeMap,
		Elem: s,
	}
}
