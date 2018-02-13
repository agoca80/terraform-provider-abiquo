# Provider configuration: Get these from the environment
variable "endpoint" { }

provider "abiquo" {
  endpoint       = "${var.endpoint}"
  username       = "terraform demo"
  password       = "12qwaszx"
}

data "abiquo_vdc"  "vdc" { name = "terraform demo" }
data "abiquo_repo" "repo" { datacenter = "datacenter 1" }

resource "abiquo_vapp" "vapp" {
  virtualdatacenter = "${data.abiquo_vdc.vdc.id}"  
  name              = "terraform demo"
}

resource "abiquo_vmt" "vmt" {
  cpu         = 1
  ram         = 64
  repo        = "${data.abiquo_repo.repo.id}"
  file        = "/path/to/demo.ova"
  name        = "terraform demo"
  description = "terraform demo"
}

resource "abiquo_vm" "server" {
  label     = "server"
  virtualappliance       = "${abiquo_vapp.vapp.id}"
  virtualmachinetemplate = "${abiquo_vmt.vmt.id}"
  variables = {
    profile = "server"
  }
  bootstrap = <<EOF
!/bin/sh
touch /server
echo root:temporal | chpasswd
exit 0
EOF
}

resource "abiquo_vm" "rss" {
  label     = "rss"
  variables = {
    profile = "rss"
  }
  virtualappliance       = "${abiquo_vapp.vapp.id}"
  virtualmachinetemplate = "${abiquo_vmt.vmt.id}"
  bootstrap = <<EOF
!/bin/sh
touch /rss
echo root:temporal | chpasswd
exit 0
EOF
}

resource "abiquo_vm" "v2v" {
  label     = "v2v"
  variables = {
    profile = "v2v"
  }
  virtualappliance       = "${abiquo_vapp.vapp.id}"
  virtualmachinetemplate = "${abiquo_vmt.vmt.id}"
  bootstrap = <<EOF
!/bin/sh
touch /v2v
echo root:temporal | chpasswd
exit 0
EOF
}
