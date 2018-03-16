package main

import (
	"fmt"
	"testing"

	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

type testHelper struct {
	kind   string
	media  string
	config string
}

func (th *testHelper) checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != th.kind {
			continue
		}
		res := core.Factory(th.media)
		href := rs.Primary.Attributes["id"]
		endpoint := core.NewLinkType(href, th.media)
		if err := core.Read(endpoint, res); err == nil {
			return fmt.Errorf("%s.test still exists: %s", th.kind, endpoint)
		}
	}
	return nil
}

func (th *testHelper) updateCase(t *testing.T) resource.TestCase {
	return resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: th.checkDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: th.config,
				Check: resource.ComposeTestCheckFunc(
					th.checkExists(),
				),
			},
		},
	}
}

func (th *testHelper) checkExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[th.kind+".test"]
		if !ok {
			return fmt.Errorf("%s.test not found", th.kind)
		}

		href := rs.Primary.Attributes["id"]
		endpoint := core.NewLinkType(href, th.media)
		return core.Read(endpoint, nil)
	}
}
