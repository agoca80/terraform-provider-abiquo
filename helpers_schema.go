package main

import (
	"fmt"
	"net"
	"net/url"
	"regexp"

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

func aggregate(e interface{}, t schema.ValueType, f schema.SchemaSetFunc) func(*schema.Schema) {
	return func(s *schema.Schema) {
		switch elem := e.(type) {
		case func(*schema.Schema):
			s.Elem = attribute(elem)
		default:
			s.Elem = elem
		}
		s.Type = t
		s.Set = f
	}
}

func set(e interface{}) func(*schema.Schema)  { return aggregate(e, schema.TypeSet, schema.HashString) }
func list(e interface{}) func(*schema.Schema) { return aggregate(e, schema.TypeList, nil) }
func hash(e interface{}) func(*schema.Schema) { return aggregate(e, schema.TypeMap, nil) }

func setFn(e interface{}, s schema.SchemaSetFunc) func(*schema.Schema) {
	return aggregate(e, schema.TypeSet, s)
}

func min(m int) func(*schema.Schema) {
	return func(s *schema.Schema) {
		s.MinItems = m
	}
}

func port(s *schema.Schema) {
	integer(s)
	s.ValidateFunc = func(d interface{}, key string) (strs []string, errs []error) {
		port := d.(int)
		if port < 1 && 65535 < port {
			errs = append(errs, fmt.Errorf("%v is an invalid port", key))
		}
		return
	}
}

func protocol(s *schema.Schema) {
	text(s)
	s.ValidateFunc = validation.StringInSlice([]string{"TCP", "HTTP", "HTTPS"}, false)
}

func email(s *schema.Schema) {
	text(s)
	s.ValidateFunc = func(d interface{}, key string) (strs []string, errs []error) {
		const emailRe = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
		if !regexp.MustCompile(emailRe).MatchString(d.(string)) {
			errs = append(errs, fmt.Errorf("%v is an invalid email", d.(string)))
		}
		return
	}
}

func price(s *schema.Schema) {
	s.Type = schema.TypeFloat
	s.Optional = true
	s.ValidateFunc = func(d interface{}, key string) (strs []string, errs []error) {
		if 0 > d.(float64) {
			errs = append(errs, fmt.Errorf("price should be 0 or greater"))
		}
		return
	}
}

func atLeast(m int) func(*schema.Schema) {
	return func(s *schema.Schema) {
		integer(s)
		s.ValidateFunc = validation.IntAtLeast(m)
	}
}

func natural(s *schema.Schema)  { atLeast(0)(s) }
func positive(s *schema.Schema) { atLeast(1)(s) }

func ip(s *schema.Schema) {
	text(s)
	s.ValidateFunc = func(d interface{}, key string) (strs []string, errs []error) {
		if net.ParseIP(d.(string)) == nil {
			errs = append(errs, fmt.Errorf("%v is an invalid IP", d.(string)))
		}
		return
	}
}

func timestamp(s *schema.Schema) {
	text(s)
	s.ValidateFunc = func(d interface{}, k string) (strs []string, errs []error) {
		if d.(string) != "" {
			strs, errs = validation.ValidateRFC3339TimeString(d, k)
		}
		return
	}
}

func href(s *schema.Schema) {
	text(s)
	s.ValidateFunc = func(d interface{}, key string) (strs []string, errs []error) {
		if _, err := url.Parse(d.(string)); err != nil {
			errs = append(errs, fmt.Errorf("%v is an invalid href", d.(string)))
		}
		return
	}
}

func link(media string) func(*schema.Schema) {
	return func(s *schema.Schema) {
		text(s)
		s.ValidateFunc = func(d interface{}, key string) (strs []string, errs []error) {
			for _, re := range validateMedia[media] {
				if regexp.MustCompile(re + "$").MatchString(d.(string)) {
					return
				}
			}
			errs = append(errs, fmt.Errorf("invalid %v : %v", key, d.(string)))
			errs = append(errs, fmt.Errorf("%v", validateMedia[media]))
			return
		}
	}
}

func label(strs []string) func(*schema.Schema) {
	return func(s *schema.Schema) {
		text(s)
		s.ValidateFunc = validation.StringInSlice(strs, false)
	}
}

func conflicts(strs []string) func(*schema.Schema) {
	return func(s *schema.Schema) {
		s.ConflictsWith = strs
	}
}

func attribute(fields ...func(*schema.Schema)) (media *schema.Schema) {
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

func byDefault(i interface{}) func(*schema.Schema) {
	return func(s *schema.Schema) {
		s.Default = i
	}
}

func variable(name string) func(*schema.Schema) {
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

func endpoint(media string) *schema.Schema {
	return attribute(func(s *schema.Schema) {
		required(s)
		forceNew(s)
		link(media)(s)
	})
}
