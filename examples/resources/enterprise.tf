resource "abiquo_enterprise" "example" {
  name = "terraform example"
  cpusoft  = 256   , cpuhard  = 512
  ramsoft  = 20480 , ramhard  = 40960
  vlansoft = 128   , vlanhard = 256
}

# PENDING: This resource should be inside the enterprise resource
resource "abiquo_limit" "example" {
  enterprise = "${abiquo_enterprise.example.id}"
  location   = "${data.abiquo_datacenter.datacenter.id}"

  cpusoft  = 250   , cpuhard  = 510
  ramsoft  = 20000 , ramhard  = 40000
  vlansoft = 120   , vlanhard = 250
  
  # PENDING: A set of available hwprofiles
  # hwprofiles = [ ]
}

resource "abiquo_role" "example" {
  name       = "terraform example"
  enterprise = "${abiquo_enterprise.example.id}"
  privileges = [
    "USERS_VIEW",
    "VDC_ENUMERATE"
  ]
}

data "abiquo_enterprise" "hwprofiles" { name = "hwprofiles" }
data "abiquo_enterprise" "neutron"    { name = "neutron" }
data "abiquo_enterprise" "oauth"      { name = "oauth" }

data "abiquo_datacenter" "dc1" { name = "datacenter 1" }
data "abiquo_datacenter" "dc2" { name = "datacenter 2" }

data "abiquo_scope" "global" { name = "Global scope" }

resource "abiquo_scope" "example" {
  name        = "terraform example"
  parent      = "${data.abiquo_scope.global.id}"
  
  datacenters = [
    "${data.abiquo_datacenter.dc1.id}",
    "${data.abiquo_datacenter.dc2.id}"
  ]
  enterprises = [
    "${data.abiquo_enterprise.hwprofiles.id}",
    "${data.abiquo_enterprise.neutron.id}",
    "${data.abiquo_enterprise.oauth.id}"
  ]
}

resource "abiquo_user" "example" {
  enterprise = "${abiquo_enterprise.example.id}"
  scope      = "${abiquo_scope.example.id}"
  role       = "${abiquo_role.example.id}"
  active     = true
  name       = "terraform example"
  surname    = "terraform example"
  nick       = "terraform example"
  email      = "terraform@abiquo.com"
}
