package main

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"time"

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
func decimal(s *schema.Schema) { s.Type = schema.TypeFloat }
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
	validate(s, func(d interface{}, key string) (strs []string, errs []error) {
		port := d.(int)
		if port < 0 && 65535 < port {
			errs = append(errs, fmt.Errorf("%v is an invalid value for %v", port, key))
		}
		return
	})
}

func protocol(s *schema.Schema) {
	text(s)
	validate(s, validation.StringInSlice([]string{"TCP", "HTTP", "HTTPS"}, false))
}

func validate(s *schema.Schema, validate schema.SchemaValidateFunc) {
	s.ValidateFunc = validate
}

func price(s *schema.Schema) {
	s.Default = 0.0
	s.Optional = true
	s.Type = schema.TypeFloat
	s.ValidateFunc = func(d interface{}, key string) (strs []string, errs []error) {
		if 0 > d.(float64) {
			errs = append(errs, fmt.Errorf("prize should be 0 or greater"))
		}
		return
	}
}

func natural(s *schema.Schema) {
	integer(s)
	validate(s, validation.IntAtLeast(0))
}

func ip(s *schema.Schema) {
	text(s)
	validate(s, func(d interface{}, key string) (strs []string, errs []error) {
		if net.ParseIP(d.(string)) == nil {
			errs = append(errs, fmt.Errorf("%v is an invalid IP", d.(string)))
		}
		return
	})
}

const tsFormat = "2006/01/02 15:04"

func timestamp(s *schema.Schema) {
	integer(s)
	s.ValidateFunc = func(d interface{}, key string) (strs []string, errs []error) {
		if _, err := time.Parse(tsFormat, d.(string)); err != nil {
			errs = append(errs, fmt.Errorf("%v is an invalid date", d.(string)))
		}
		return
	}
}

func href(s *schema.Schema) {
	text(s)
	validate(s, func(d interface{}, key string) (strs []string, errs []error) {
		if _, err := url.Parse(d.(string)); err != nil {
			errs = append(errs, fmt.Errorf("%v is an invalid IP", d.(string)))
		}
		return
	})
}

func link(s *schema.Schema, regexps []string) {
	text(s)
	validate(s, func(d interface{}, key string) (strs []string, errs []error) {
		for _, re := range regexps {
			if regexp.MustCompile(re).MatchString(d.(string)) {
				return
			}
		}
		errs = append(errs, fmt.Errorf("invalid %v value : %v", key, d.(string)))
		return
	})
}

func label(strs []string) field {
	return func(s *schema.Schema) {
		text(s)
		validate(s, validation.StringInSlice(strs, false))
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
