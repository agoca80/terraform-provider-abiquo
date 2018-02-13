data "abiquo_datacenter" "datacenter" { name = "datacenter 1" }

resource "abiquo_hp" "small" {
  active = true
  name = "terraform small"
  cpu  = 1
  ram  = 64
  datacenter = "${data.abiquo_datacenter.datacenter.id}"
}

resource "abiquo_hp" "medium" {
  active = true
  name = "terraform medium"
  cpu  = 1
  ram  = 512
  datacenter = "${data.abiquo_datacenter.datacenter.id}"
}

resource "abiquo_hp" "big" {
  active = true
  name = "terraform big"
  cpu  = 1
  ram  = 1024
  datacenter = "${data.abiquo_datacenter.datacenter.id}"
}

data "abiquo_enterprise" "enterprise" {
  name = "hwprofiles"
}

# PENDING: This resource should be inside the enterprise resource
resource "abiquo_limit" "example" {
  enterprise = "${data.abiquo_enterprise.enterprise.id}"
  location   = "${data.abiquo_datacenter.datacenter.id}"
  
  # PENDING: A set of available hwprofiles
  hwprofiles = [ 
    "${abiquo_hp.small.id}",
    "${abiquo_hp.medium.id}",
    "${abiquo_hp.big.id}"
  ]
}

