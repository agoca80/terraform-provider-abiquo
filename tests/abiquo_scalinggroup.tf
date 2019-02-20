resource "abiquo_scalinggroup" "test" {
  mastervirtualmachine = "${abiquo_virtualmachine.test.id}"
  virtualappliance     = "${abiquo_virtualappliance.test.id}"

  name      = "test sg"
  cooldown  = 60
  min       = 0
  max       = 4
  scale_in  = [ { numberofinstances = 1 } ]
  scale_out = [ { numberofinstances = 1 } ]
}

data "abiquo_virtualdatacenter"      "test"     { name = "${var.virtualdatacenter}" }
data "abiquo_template" "test"     {
  templates = "${data.abiquo_virtualdatacenter.test.templates}"
  name = "${var.template}"
}

resource "abiquo_virtualappliance" "test" {
  virtualdatacenter = "${data.abiquo_virtualdatacenter.test.id}"
  name = "test"
}

# Scaling group master instance
resource "abiquo_virtualmachine" "test" {
  deploy                 = false
  virtualappliance       = "${abiquo_virtualappliance.test.id}"
  virtualmachinetemplate = "${data.abiquo_template.test.id}"
  label                  = "test sg"
}
