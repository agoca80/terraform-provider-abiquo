resource "abiquo_plan" "test" {
	virtualmachine = "${abiquo_vm.test.id}"
	description    = "test plan"
	name           = "test plan"
	entries        = [
		{	parameter = "", parametertype = "None",	type = "UNDEPLOY" }
	]
}

data "abiquo_vdc"        "test"       { name = "tests" }
data "abiquo_template"   "test"       {
  templates = "${data.abiquo_vdc.test.templates}"
  name      = "tests"
}

resource "abiquo_vapp" "test" {
	virtualdatacenter = "${data.abiquo_vdc.test.id}"
	name              = "test plan"
}

resource "abiquo_vm" "test" {
	cpu                    = 1
	deploy                 = false
	ram                    = 64
	label                  = "test plan"
	virtualappliance       = "${abiquo_vapp.test.id}"
	virtualmachinetemplate = "${data.abiquo_template.test.id}"
}
