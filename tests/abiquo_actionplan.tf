resource "abiquo_actionplan" "test" {
	description    = "test plan"
	name           = "test plan"
	entries        = [
		{
			parameter = ""
			parametertype = "None"
			type = "UNDEPLOY"
			links = [
				"${abiquo_virtualmachine.test.id}"
			]
		}
	]
}

data "abiquo_virtualdatacenter"        "test"       { name = "tests" }
data "abiquo_template"   "test"       {
  templates = "${data.abiquo_virtualdatacenter.test.templates}"
  name      = "tests"
}

resource "abiquo_virtualappliance" "test" {
	virtualdatacenter = "${data.abiquo_virtualdatacenter.test.id}"
	name              = "test plan"
}

resource "abiquo_virtualmachine" "test" {
	cpu                    = 1
	deploy                 = false
	ram                    = 64
	label                  = "test plan"
	virtualappliance       = "${abiquo_virtualappliance.test.id}"
	virtualmachinetemplate = "${data.abiquo_template.test.id}"
}
