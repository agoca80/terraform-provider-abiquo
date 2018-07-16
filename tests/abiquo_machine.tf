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

data "abiquo_dstier" "test" {
  datacenter = "${data.abiquo_datacenter.test.id}"
  name       = "Default Tier"
}

variable "kvm" {}
variable "iface" {}
variable "datastore" {}

data "abiquo_machine" "test" {
  datacenter = "${data.abiquo_datacenter.test.id}"
  hypervisor = "KVM"
  ip         = "${var.kvm}"
}

resource "abiquo_machine" "test" {
  rack       = "${abiquo_rack.test.id}"
  definition = "${data.abiquo_machine.test.definition}"

  interface {
    name = "${var.iface}"
    nst  = "${data.abiquo_nst.test.id}"
  }

  datastore {
    uuid   = "${var.datastore}"
    dstier = "${data.abiquo_dstier.test.id}"
  }

  lifecycle = {
    "ignore_changes" = ["definition"]
  }
}
