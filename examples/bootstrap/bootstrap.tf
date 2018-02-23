# Provider configuration: Get these from the environment
variable "endpoint" { }
variable "username" { }
variable "password" { }

provider "abiquo" {
  endpoint       = "${var.endpoint}"
  username       = "${var.username}"
  password       = "${var.password}"
}

variable "am"         { }
variable "rss"        { }
variable "license"    { }
variable "dcname"     { }
variable "dclocation" { }

resource "abiquo_license" "license" {
  code = "${var.license}"
}

resource "abiquo_datacenter" "datacenter" {
  name     = "${var.dcname}"
  location = "${var.dclocation}"
  am       = "${var.am}"
  vf       = "http://${var.rss}:8009/virtualfactory"
  nc       = "http://${var.rss}:8009/nodecollector"
  bpm      = "http://${var.rss}:8009/bpm-async"
  ssm      = "http://${var.rss}:8009/ssm"
  vsm      = "http://${var.rss}:8009/vsm"
  cpp      = "http://${var.rss}:8009/cpp"
  dhcp     = "omapi://${var.rss}:7911/"
  dhcpv6   = "omapi://${var.rss}:7911/"
  ra       = "guacd://${var.rss}:4822/"
}

resource "abiquo_rack" "kvm" {
  name        = "kvm"
  vlanmin     = 1000
  vlanmax     = 1999
  description = "kvm"
  datacenter  = "${abiquo_datacenter.datacenter.id}"
}

resource "abiquo_rack" "vmx" {
  name       = "vmx"
  vlanmin     = 2000
  vlanmax     = 2999
  description = "vmx"
  datacenter = "${abiquo_datacenter.datacenter.id}"
}

data "abiquo_nst" "default" {
  datacenter = "${abiquo_datacenter.datacenter.id}"
  name       = "Service Network"
}

data "abiquo_dstier" "default" {
  datacenter = "${abiquo_datacenter.datacenter.id}"
  name       = "Default Tier"
}

data "abiquo_machine" "kvm" {
  datacenter = "${abiquo_datacenter.datacenter.id}"
  hypervisor = "KVM"
  ip         = "10.60.13.5"
}

resource "abiquo_machine" "kvm" {
  definition = "${data.abiquo_machine.kvm.definition}"
  rack       = "${abiquo_rack.kvm.id}"
  datastores = { "12452353-dd5f-4ec0-b739-cac38b54152f" = "${data.abiquo_dstier.default.id}" }
  interfaces = { "52:54:00:eb:c3:4b" = "${data.abiquo_nst.default.id}" }
  lifecycle  = { "ignore_changes" = ["definition"] }
}

variable "vcenter_address"  { }
variable "vcenter_username" { }
variable "vcenter_password" { }

data "abiquo_machine" "vmx0" {
  datacenter  = "${abiquo_datacenter.datacenter.id}"
  hypervisor  = "VMX_04"
  ip          = "192.168.2.64"
  managerip   = "${var.vcenter_address}"
  manageruser = "${var.vcenter_username}"
  managerpass = "${var.vcenter_password}"
}

data "abiquo_machine" "vmx1" {
  datacenter  = "${abiquo_datacenter.datacenter.id}"
  hypervisor  = "VMX_04"
  ip          = "192.168.2.65"
  managerip   = "${var.vcenter_address}"
  manageruser = "${var.vcenter_username}"
  managerpass = "${var.vcenter_password}"
}

resource "abiquo_machine" "vmx0" {
  definition  = "${data.abiquo_machine.vmx0.definition}"
  rack        = "${abiquo_rack.vmx.id}"
  managerip   = "${var.vcenter_address}"
  manageruser = "${var.vcenter_username}"
  managerpass = "${var.vcenter_password}"
  datastores  = { "dbd17fde-7c11-42e7-863c-5b45455577db" = "${data.abiquo_dstier.default.id}" }
  interfaces  = { "00:15:c5:ff:1c:d9" = "${data.abiquo_nst.default.id}" }
  lifecycle   = { "ignore_changes" = ["definition"] }
}

resource "abiquo_machine" "vmx1" {
  definition  = "${data.abiquo_machine.vmx1.definition}"
  rack        = "${abiquo_rack.vmx.id}"
  managerip   = "${var.vcenter_address}"
  manageruser = "${var.vcenter_username}"
  managerpass = "${var.vcenter_password}"
  datastores  = { "a745640a-2b87-4821-be3a-7e5a73ad0a96" = "${data.abiquo_dstier.default.id}" }
  interfaces  = { "00:15:c5:ff:24:72" = "${data.abiquo_nst.default.id}" }
  lifecycle   = { "ignore_changes" = ["definition"] }
}
