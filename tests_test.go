package main

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func checkDestroy(d *description) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != d.Name() {
				continue
			}
			href := rs.Primary.Attributes["id"]
			endpoint := linkType(href, d.media)
			if err := endpoint.Read(nil); err == nil {
				return fmt.Errorf("%s.test still exists: %s", d.Name(), endpoint.URL())
			}
		}
		return nil
	}
}

func checkExists(d *description) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[d.Name()+".test"]
		if !ok {
			return fmt.Errorf("%s.test not found", d.Name())
		}

		href := rs.Primary.Attributes["id"]
		endpoint := linkType(href, d.media)
		_, err := endpoint.Exists()
		return err
	}
}

func basicTest(t *testing.T, d *description) {
	file := "tests/" + d.Name() + ".tf"
	config, err := ioutil.ReadFile(file)
	if err != nil {
		t.Error("updateCase:", file, "could not be read:", err)
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkDestroy(d),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: string(config),
				Check: resource.ComposeTestCheckFunc(
					checkExists(d),
				),
			},
		},
	})
}

func TestActionPlan(t *testing.T)      { basicTest(t, actionplan) }
func TestAlarm(t *testing.T)           { basicTest(t, alarm) }
func TestAlert(t *testing.T)           { basicTest(t, alert) }
func TestBackuppolicy(t *testing.T)    { basicTest(t, backuppolicy) }
func TestCostcode(t *testing.T)        { basicTest(t, costcode) }
func TestCurrency(t *testing.T)        { basicTest(t, currency) }
func TestDatastorerule(t *testing.T)   { basicTest(t, datastoreloadrule) }
func TestDatastoretier(t *testing.T)   { basicTest(t, datastoretier) }
func TestDevice(t *testing.T)          { basicTest(t, device) }
func TestEnterprise(t *testing.T)      { basicTest(t, enterprise) }
func TestExternal(t *testing.T)        { basicTest(t, external) }
func TestFirewallpolicy(t *testing.T)  { basicTest(t, firewallpolicy) }
func TestHardwareprofile(t *testing.T) { basicTest(t, hardwareprofile) }
func TestLoadbalancer(t *testing.T)    { basicTest(t, loadbalancer) }
func TestLimit(t *testing.T)           { basicTest(t, limit) }
func TestMachine(t *testing.T)         { basicTest(t, machine) }
func TestMachinerule(t *testing.T)     { basicTest(t, machineloadrule) }
func TestPricing(t *testing.T)         { basicTest(t, pricingtemplate) }
func TestPrivate(t *testing.T)         { basicTest(t, private) }
func TestPublic(t *testing.T)          { basicTest(t, public) }
func TestRack(t *testing.T)            { basicTest(t, rack) }
func TestRole(t *testing.T)            { basicTest(t, role) }
func TestScope(t *testing.T)           { basicTest(t, scope) }
func TestScalinggroup(t *testing.T)    { basicTest(t, scalinggroup) }
func TestTemplate(t *testing.T)        { basicTest(t, virtualmachinetemplate) }
func TestUser(t *testing.T)            { basicTest(t, user) }
func TestVAPP(t *testing.T)            { basicTest(t, virtualappliance) }
func TestVDC(t *testing.T)             { basicTest(t, virtualdatacenter) }
func TestVM(t *testing.T)              { basicTest(t, virtualmachine) }
func TestVolume(t *testing.T)          { basicTest(t, volume) }
