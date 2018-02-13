data "abiquo_repo" "repo" {
  datacenter = "datacenter 1"
}

resource "abiquo_vmt" "vmt" {
  cpu         = 1
  ram         = 64
  repo        = "${data.abiquo_repo.repo.id}"
  file        = "/path/to/test.ova"
  name        = "terraform examples"
  description = "A template uploaded from terraform"
}
