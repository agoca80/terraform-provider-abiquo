package main

import "github.com/hashicorp/terraform/helper/schema"

type description schema.Schema

func Optional() *description {
	return &description{
		Optional: true,
	}
}

func Required() *description {
	return &description{
		Required: true,
	}
}

func Computed() *description {
	return &description{
		Computed: true,
	}
}

func Conflicts(conflicts []string) *description {
	return &description{
		ConflictsWith: conflicts,
		Optional:      true,
	}
}

func (d *description) set(s func(interface{}) int, min int, kind schema.ValueType) *schema.Schema {
	d.Set = s
	d.MinItems = min
	d.Type = schema.TypeSet
	d.Elem = &schema.Schema{
		Type: kind,
	}
	return d.schema()
}

func (d *description) list(min int, elem interface{}) *schema.Schema {
	d.Type = schema.TypeList
	d.MinItems = min
	d.Elem = elem
	return d.schema()
}

func (d *description) Map(kind schema.ValueType) *schema.Schema {
	d.Type = schema.TypeMap
	d.Elem = &schema.Schema{
		Type: kind,
	}
	return d.schema()
}

func (d *description) Renew() *description {
	d.ForceNew = true
	return d
}

func (d *description) Validate(f schema.SchemaValidateFunc) *description {
	d.ValidateFunc = f
	return d
}

func (d *description) schema() *schema.Schema {
	s := schema.Schema(*d)
	return &s
}

func (d *description) String() *schema.Schema {
	d.Type = schema.TypeString
	return d.schema()
}

func (d *description) Number() *schema.Schema {
	d.Type = schema.TypeInt
	return d.schema()
}

func (d *description) Bool() *schema.Schema {
	d.Type = schema.TypeBool
	return d.schema()
}

func (d *description) Link() *schema.Schema {
	d.Type = schema.TypeString
	d.ValidateFunc = validateURL
	return d.schema()
}

func (d *description) Timestamp() *schema.Schema {
	d.Type = schema.TypeString
	d.ValidateFunc = validateTS
	return d.schema()
}

func (d *description) Links() *schema.Schema {
	d.Type = schema.TypeList
	d.Elem = &schema.Schema{
		Type:         schema.TypeString,
		ValidateFunc: validateURL,
	}
	return d.schema()
}

func (d *description) IP() *schema.Schema {
	return d.Validate(validateIP).String()
}

func (d *description) ValidateString(strings []string) *schema.Schema {
	d.Type = schema.TypeString
	d.ValidateFunc = validateString(strings)
	return d.schema()
}

func (d *description) ValidateURL() *schema.Schema {
	d.Type = schema.TypeString
	d.ValidateFunc = validateURL
	return d.schema()
}
