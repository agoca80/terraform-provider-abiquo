package main

import (
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

var (
	validateProtocol   = validateString([]string{"ALL", "TCP", "UDP"})
	validateAlgorithm  = validateString([]string{"ROUND_ROBIN", "LEAST_CONNECTIONS", "SOURCE_IP"})
	validateController = validateString([]string{"IDE", "SCSI", "VIRTIO"})
)

func validateString(values []string) schema.SchemaValidateFunc {
	return func(d interface{}, key string) (strs []string, errs []error) {
		for _, v := range values {
			if v == d {
				return
			}
		}
		errs = append(errs, fmt.Errorf("%s is an invalid value for %s", d.(string), key))
		return
	}
}

func validateIP(d interface{}, key string) (strs []string, errs []error) {
	if net.ParseIP(d.(string)) == nil {
		errs = append(errs, fmt.Errorf("%v is an invalid IP", d.(string)))
	}
	return
}

func validateURL(d interface{}, key string) (strs []string, errs []error) {
	if _, err := url.Parse(d.(string)); err != nil {
		errs = append(errs, fmt.Errorf("%v is an invalid IP", d.(string)))
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
