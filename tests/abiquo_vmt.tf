resource "abiquo_virtualmachinetemplate" "test" {
  cpu         = 1
  ram         = 64
  repo        = "${data.abiquo_repo.repo.id}"
  file        = "${var.test_ova}"
  name        = "test virtualmachinetemplate"
  description = "test virtualmachinetemplate"
}

variable "test_ova" {  }
data     "abiquo_repo" "repo" { datacenter = "datacenter 1" }
