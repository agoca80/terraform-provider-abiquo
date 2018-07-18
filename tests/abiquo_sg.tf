resource "abiquo_sg" "test" {
  mastervirtualmachine = "${abiquo_vm.test.id}"
  virtualappliance     = "${abiquo_vapp.test.id}"

  name      = "test sg"
  cooldown  = 60
  min       = 0
  max       = 4
  scale_in  = [ { numberofinstances = 1 } ]
  scale_out = [ { numberofinstances = 1 } ]
}

data "abiquo_vdc"      "test"     { name = "tests" }
data "abiquo_template" "test"     {
  templates = "${data.abiquo_vdc.test.templates}"
  name      = "tests"
}

resource "abiquo_vapp" "test" {
  virtualdatacenter = "${data.abiquo_vdc.test.id}"
  name = "test"
}

# Scaling group master instance
resource "abiquo_vm" "test" {
  deploy                 = false
  virtualappliance       = "${abiquo_vapp.test.id}"
  virtualmachinetemplate = "${data.abiquo_template.test.id}"
  label                  = "test sg"
}
