resource "abiquo_machine" "test" {
  rack       = "${abiquo_rack.test.id}"
  definition = "${data.abiquo_machine.test.definition}"

  interface {
    name = "${var.test_kvm_interface}"
    nst  = "${data.abiquo_nst.test.id}"
  }

  datastore {
    uuid   = "${var.test_kvm_datastore}"
    dstier = "${data.abiquo_datastoretier.test.id}"
  }

  lifecycle = {
    "ignore_changes" = ["definition"]
  }
}

data "abiquo_datacenter" "test" {
  name = "datacenter 1"
}

resource "abiquo_rack" "test" {
  name        = "test machine"
  description = "kvm"
  datacenter  = "${data.abiquo_datacenter.test.id}"
}

data "abiquo_nst" "test" {
  datacenter = "${data.abiquo_datacenter.test.id}"
  name       = "Service Network"
}

data "abiquo_datastoretier" "test" {
  datacenter = "${data.abiquo_datacenter.test.id}"
  name       = "Default Tier"
}

variable "test_kvm_ip" {}
variable "test_kvm_interface" {}
variable "test_kvm_datastore" {}

data "abiquo_machine" "test" {
  datacenter = "${data.abiquo_datacenter.test.id}"
  hypervisor = "KVM"
  ip         = "${var.test_kvm_ip}"
}
