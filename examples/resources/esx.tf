# This example requires a

resource "abiquo_vdc" "esx" {
  enterprise = "${data.abiquo_enterprise.enterprise.id}"
  location   = "${data.abiquo_location.location.id}"

  name       = "terraform"
  type       = "VMX_04"
}

resource "abiquo_private" "esx" {
  count = 3
  
  virtualdatacenter = "${abiquo_vdc.esx.id}"

  # XXX workaround ABICLOUDPREMIUM-9660
  lifecycle = { ignore_changes = [ "dns1", "dns2" ] }

  mask    = 24
  address = "172.16.${count.index}.0"
  gateway = "172.16.${count.index}.1"
  name    = "terraform ${count.index}"
  dns1    = "8.8.8.8"
  dns2    = "4.4.4.4"
  suffix  = "private${count.index}.bcn.com"
}

resource "abiquo_vapp" "esx" {
  virtualdatacenter = "${abiquo_vdc.esx.id}"  
  name = "terraform"
}

resource "abiquo_vol" "esx" {
  tier               = "${abiquo_vdc.esx.id}/tiers/4"
  virtualdatacenter  = "${abiquo_vdc.esx.id}"

  type = "SCSI"
  name = "terraform"
  size = 32
}

resource "abiquo_hd" "esx" {
  virtualdatacenter  = "${abiquo_vdc.esx.id}"
  type  = "SCSI"
  size  = 32
  label = "terraform"
}

variable "bootstrap" {
  type = "string"
  default = <<EOF
!/bin/sh
touch /bootstrapped
echo root:temporal | chpasswd
exit 0
EOF
}

resource "abiquo_vm" "esx" {
  virtualappliance       = "${abiquo_vapp.esx.id}"
  virtualmachinetemplate = "${data.abiquo_template.template.id}"

  label     = "terraform"
  bootstrap = "${var.bootstrap}"

  disks = [
     "${abiquo_hd.esx.id}",
     "${abiquo_vol.esx.id}"
  ]
  
  ips = [
    "${abiquo_private.esx.0.id}/ips",
    "${abiquo_private.esx.1.id}/ips",
    "${abiquo_private.esx.2.id}/ips"
  ]  
}
