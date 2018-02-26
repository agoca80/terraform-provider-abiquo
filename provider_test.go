package main

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"abiquo": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	endpoint := os.Getenv("ABQ_ENDPOINT")
	username := os.Getenv("ABQ_USERNAME")
	password := os.Getenv("ABQ_PASSWORD")
	if endpoint == "" || username == "" || password == "" {
		t.Fatal("ABQ_ENDPOINT, ABQ_USERNAME and ABQ_PASSWORD must be set for acceptance tests")
	}
}
