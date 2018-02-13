# This example requires a

resource "abiquo_vdc" "deploy" {
  enterprise = "${data.abiquo_enterprise.enterprise.id}"
  location   = "${data.abiquo_location.location.id}"

  name       = "terraform deploy"
  type       = "VMX_04"
}

resource "abiquo_vapp" "deploy" {
  virtualdatacenter = "${abiquo_vdc.deploy.id}"  
  name = "terraform deploy"
}

data "abiquo_template" "template" { name = "tests" }

resource "abiquo_vm" "deploy" {
  label     = "terraform deploy"
  variables = {
    variable1 = "terraform variable1"
    variable2 = "terraform variable2"
  }
  virtualappliance       = "${abiquo_vapp.deploy.id}"
  virtualmachinetemplate = "${data.abiquo_template.template.id}"
  bootstrap = <<EOF
#!/bin/sh
exit 0
EOF
}
