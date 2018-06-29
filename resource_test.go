package main

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func checkDestroy(name, media string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != name {
				continue
			}
			href := rs.Primary.Attributes["id"]
			endpoint := core.NewLinkType(href, media)
			if err := core.Read(endpoint, nil); err == nil {
				return fmt.Errorf("%s.test still exists: %s", name, endpoint)
			}
		}
		return nil
	}
}

func checkExists(name, media string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name+".test"]
		if !ok {
			return fmt.Errorf("%s.test not found", name)
		}

		href := rs.Primary.Attributes["id"]
		endpoint := core.NewLinkType(href, media)
		return core.Read(endpoint, nil)
	}
}

func updateCase(t *testing.T, name, media string) resource.TestCase {
	file := "tests/" + name + ".tf"
	config, err := ioutil.ReadFile(file)
	if err != nil {
		t.Error("updateCase:", file, "could not be read:", err)
	}
	return resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkDestroy(name, media),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: string(config),
				Check: resource.ComposeTestCheckFunc(
					checkExists(name, media),
				),
			},
		},
	}
}

func TestAlarm_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_alarm", "alarm"))
}

func TestAlert_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_alert", "alert"))
}

func TestBackup_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_backup", "backuppolicy"))
}

func TestCompute_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_computeload", "machineloadrule"))
}

func TestCostCode_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_costcode", "costcode"))
}

func TestCurrency_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_currency", "currency"))
}

func TestDatastoreTier_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_dstier", "datastoretier"))
}

func TestDevice_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_device", "device"))
}

func TestEnterprise_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_enterprise", "enterprise"))
}

func TestAccExternal_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_external", "vlan"))
}

func TestFW_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_fw", "firewallpolicy"))
}

func TestHP_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_hp", "hardwareprofile"))
}

func TestLB_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_lb", "loadbalancer"))
}

func TestLimit_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_limit", "limit"))
}

func TestPlan_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_plan", "virtualmachineactionplan"))
}

func TestPricing_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_pricing", "pricingtemplate"))
}

func TestPublic_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_public", "vlan"))
}

func TestRack_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_rack", "rack"))
}

func TestRole_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_role", "role"))
}

func TestScope_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_scope", "scope"))
}

func TestSG_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_sg", "scalinggroup"))
}

func TestStorageLoad_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_storageload", "datastoreloadrule"))
}

func TestUser_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_user", "user"))
}

func TestVapp_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_vapp", "virtualappliance"))
}

func TestVdc_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_vdc", "virtualdatacenter"))
}

func TestVM_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_vm", "virtualmachine"))
}

func TestVMT_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_vmt", "virtualmachinetemplate"))
}

func TestVolume_update(t *testing.T) {
	resource.Test(t, updateCase(t, "abiquo_vol", "volume"))
}
