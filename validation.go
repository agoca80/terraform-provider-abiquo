package main

import (
	"fmt"
	"net"
	"net/url"
	"time"
)

func validatePort(d interface{}, key string) (strs []string, errs []error) {
	port := d.(int)
	if port < 0 && 65535 < port {
		errs = append(errs, fmt.Errorf("%v is an invalid value for %v", port, key))
	}
	return
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

func validatePrice(d interface{}, key string) (strs []string, errs []error) {
	if 0 > d.(float64) {
		errs = append(errs, fmt.Errorf("prize should be 0 or greater"))
	}
	return
}
