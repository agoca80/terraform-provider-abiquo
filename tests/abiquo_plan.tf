data "abiquo_template"   "template"   { name = "tests" }
data "abiquo_vdc"        "test"       { name = "tests" }

resource "abiquo_vapp" "vapp" {
	virtualdatacenter = "${data.abiquo_vdc.test.id}"
	name              = "test plan"
}

resource "abiquo_vm" "vm" {
	cpu                    = 1
	deploy                 = false
	ram                    = 64
	label                  = "test plan"
	virtualappliance       = "${abiquo_vapp.vapp.id}"
	virtualmachinetemplate = "${data.abiquo_template.template.id}"
}

resource "abiquo_plan" "test" {
	virtualmachine = "${abiquo_vm.vm.id}"
	description    = "test plan"
	name           = "test plan"
	entries        = [
		{	parameter = "", parametertype = "None",	type = "UNDEPLOY" }
	]
}
