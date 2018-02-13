# Provider configuration: Get these from the environment
variable "endpoint" { }
variable "username" { }
variable "password" { }

provider "abiquo" {
  endpoint       = "${var.endpoint}"
  username       = "${var.username}"
  password       = "${var.password}"
}

data "abiquo_enterprise" "enterprise" { name = "terraform demo" }
data "abiquo_datacenter" "datacenter" { name = "datacenter 1" }

resource "abiquo_scope" "scope" {
  name        = "terraform demo"
  datacenters = [ "${data.abiquo_datacenter.datacenter.id}" ]
  enterprises = [ "${data.abiquo_enterprise.enterprise.id}" ]
}

resource "abiquo_role" "role" {
  name = "terraform demo"
  privileges = [
    "APPLIB_UPLOAD_IMAGE",
    "APPLIB_ALLOW_MODIFY",
    "APPLIB_VIEW",
    "VAPP_CREATE_STATEFUL",
    "VAPP_CUSTOMISE_SETTINGS",
    "VAPP_DEPLOY_UNDEPLOY",
    "VAPP_RESTORE_BACKUP",
    "VAPP_STATEFUL_VIEW",
    "VDC_ENUMERATE",
    "VDC_MANAGE_STORAGE",
    "VDC_MANAGE_STORAGE_CONTROLLER",
    "VDC_MANAGE_STORAGE_DISK_ALLOCATION",
    "VDC_MANAGE_VAPP",
    "VM_ACTION_PLAN_MANAGE",
    "VM_ATTACH_NIC",
    "VM_CHECK_USER_PASSWORD",
    "VM_EDIT_CPU_RAM",
    "VM_PROTECT_ACTION",
  ]
}

resource "abiquo_user" "user" {
  enterprise = "${data.abiquo_enterprise.enterprise.id}"
  scope      = "${abiquo_scope.scope.id}"
  role       = "${abiquo_role.role.id}"
  active     = true
  name       = "terraform demo"
  surname    = "terraform demo"
  nick       = "terraform demo"
  email      = "example@example.com"
}

data "abiquo_location"   "datacenter" { name = "datacenter 1" }

resource "abiquo_vdc" "vdc" {
  enterprise = "${data.abiquo_enterprise.enterprise.id}"
  location   = "${data.abiquo_location.datacenter.id}"

  name       = "terraform demo"
  type       = "VMX_04"
}

# resource "abiquo_external" "external" {
#   enterprise         = "${data.abiquo_enterprise.enterprise.id}"
#   datacenter         = "${data.abiquo_datacenter.datacenter.id}"
#   networkservicetype = "${data.abiquo_datacenter.datacenter.id}/networkservicetypes/1"
# 
#   # XXX workaround ABICLOUDPREMIUM-9660
#   lifecycle = { ignore_changes = [ "dns1", "dns2" ] }
# 
#   tag     = 1331
#   mask    = 24
#   address = "172.16.4.0"
#   gateway = "172.16.4.1"
#   name    = "terraform demo"
#   dns1    = "4.4.4.4"
#   dns2    = "8.8.8.8"
# }
