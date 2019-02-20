resource "abiquo_virtualmachinetemplate" "test" {
  cpu         = 1
  ram         = 64
  repo        = "${data.abiquo_repo.repo.id}"
  ova         = "${var.ova}"
  name        = "test virtualmachinetemplate"
  description = "test virtualmachinetemplate"
}

variable "ova" {  }

data     "abiquo_repo" "repo" { datacenter = "${var.datacenter}" }
