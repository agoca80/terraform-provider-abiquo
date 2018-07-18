resource "abiquo_vmt" "test" {
  cpu         = 1
  ram         = 64
  repo        = "${data.abiquo_repo.repo.id}"
  file        = "${var.ova}"
  name        = "test vmt"
  description = "test vmt"
}

variable "ova" {  }
data     "abiquo_repo" "repo" { datacenter = "datacenter 1" }
