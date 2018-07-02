package main

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

var validateMedia = map[string]schema.SchemaValidateFunc{
	// admin
	"backuppolicy_dc":    validateLink([]string{"/admin/datacenters/[0-9]+/backuppolicies/[0-9]+"}),
	"datacenter":         validateLink([]string{"/admin/datacenters/[0-9]+"}),
	"datastoretier_dc":   validateLink([]string{"/admin/datacenters/[0-9]+/datastoretiers/[0-9]+"}),
	"enterprise":         validateLink([]string{"/admin/enterprises/[0-9]+"}),
	"hardwareprofile_dc": validateLink([]string{"/admin/datacenters/[0-9]+/hardwareprofiles/[0-9]+"}),
	"storagedevice":      validateLink([]string{"/admin/datacenters/[0-9]+/storage/devices/[0-9]+"}),
	// cloud
	"location":          validateLink([]string{"/cloud/locations/[0-9]+"}),
	"backuppolicy_vdc":  validateLink([]string{"/cloud/locations/[0-9]+/backuppolicies/[0-9]+"}),
	"privatenetwork":    validateLink([]string{"/cloud/virtualdatacenters/[0-9]+/privatenetworks/[0-9]+"}),
	"vdcTier":           validateLink([]string{"/cloud/virtualdatacenters/[0-9]+/tiers/[0-9]+"}),
	"virtualappliance":  validateLink([]string{"/cloud/virtualdatacenters/[0-9]+/virtualappliances/[0-9]+"}),
	"virtualdatacenter": validateLink([]string{"/cloud/virtualdatacenters/[0-9]+"}),
	"virtualmachine":    validateLink([]string{"/cloud/virtualdatacenters/[0-9]+/virtualappliances/[0-9]/virtualmachines/[0-9]+"}),
	"virtualmachine_ip": validateLink([]string{
		"/admin/enterprises/[0-9]+/limits/[0-9]+/externalnetworks/[0-9]+/ips/[0-9]+",
		"/cloud/virtualdatacenters/[0-9]+/privatenetworks/[0-9]+/ips/[0-9]+",
		"/cloud/virtualdatacenters/[0-9]+/publicips/purchased/[0-9]+",
	}),
}

func validateIP(d interface{}, key string) (strs []string, errs []error) {
	if net.ParseIP(d.(string)) == nil {
		errs = append(errs, fmt.Errorf("%v is an invalid IP", d.(string)))
	}
	return
}

func validatePort(d interface{}, key string) (strs []string, errs []error) {
	port := d.(int)
	if port < 1 && 65535 < port {
		errs = append(errs, fmt.Errorf("%v is an invalid port", key))
	}
	return
}

func validatePrice(d interface{}, key string) (strs []string, errs []error) {
	if 0 > d.(float64) {
		errs = append(errs, fmt.Errorf("price should be 0 or greater"))
	}
	return
}

const tsFormat = "2006/01/02 15:04"

func validateTS(d interface{}, key string) (strs []string, errs []error) {
	if _, err := time.Parse(tsFormat, d.(string)); err != nil {
		errs = append(errs, fmt.Errorf("%v is an invalid date", d.(string)))
	}
	return
}

func validateHref(d interface{}, key string) (strs []string, errs []error) {
	if _, err := url.Parse(d.(string)); err != nil {
		errs = append(errs, fmt.Errorf("%v is an invalid href", d.(string)))
	}
	return
}

func validateLink(regexps []string) schema.SchemaValidateFunc {
	return func(d interface{}, key string) (strs []string, errs []error) {
		for _, re := range regexps {
			if regexp.MustCompile(re + "$").MatchString(d.(string)) {
				return
			}
		}
		errs = append(errs, fmt.Errorf("invalid %v : %v", key, d.(string)))
		return
	}
}
