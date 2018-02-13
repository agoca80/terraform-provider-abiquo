data "abiquo_enterprise" "enterprise" { name = "hwprofiles" }
data "abiquo_location"   "location"   { name = "datacenter 1" }
data "abiquo_template"   "template"   { name = "tests" }
data "abiquo_hp"         "small"      {
  name     = "terraform small"
  location = "datacenter 1"
}

resource "abiquo_vdc" "vdc" {
  enterprise = "${data.abiquo_enterprise.enterprise.id}"
  location   = "${data.abiquo_location.location.id}"

  name       = "terraform example"
  type       = "VMX_04"
}

resource "abiquo_vapp" "vapp" {
  virtualdatacenter = "${abiquo_vdc.vdc.id}"  
  name = "terraform example"
}

resource "abiquo_vm" "vm" {
   virtualappliance       = "${abiquo_vapp.vapp.id}"
   virtualmachinetemplate = "${data.abiquo_template.template.id}"
   hardwareprofile        = "${data.abiquo_hp.small.id}"

   label = "terraform example"
}
