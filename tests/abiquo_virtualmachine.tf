resource "abiquo_virtualmachine" "test" {
  virtualappliance       = "${abiquo_virtualappliance.test.id}"
  virtualmachinetemplate = "${data.abiquo_template.test.id}"

  deploy  = false
  cpu     = 1
  ram     = 64
  label   = "test vm"
  backups = [ "${data.abiquo_backup.test.id}" ]
  lbs     = [ "${abiquo_loadbalancer.test.id}" ]
  fws     = [ "${abiquo_firewallpolicy.test.id}" ]

  variables = {
    name1 = "value1"
    name2 = "value2"
  }

  ips = [
    "${abiquo_ip.private.id}",
    "${data.abiquo_ip.external.id}",
    "${data.abiquo_ip.public.id}"
  ]

  bootstrap = <<EOF
#!/bin/sh
exit 0
EOF
}

data "abiquo_enterprise" "test" { name = "${var.enterprise}" }
data "abiquo_location"   "test" { name = "${var.datacenter}" }
data "abiquo_datacenter" "test" { name = "${var.datacenter}" }

data "abiquo_nst"        "test"        {
  datacenter = "${data.abiquo_datacenter.test.id}"
  name       = "${var.nst}"
}

data "abiquo_template" "test" {
  templates = "${abiquo_virtualdatacenter.test.templates}"
  name = "${var.template}"
}

resource "abiquo_backuppolicy" "test" {
  datacenter     = "${data.abiquo_datacenter.test.id}"
  code           = "test vm"
  name           = "test vm"
  description    = "test vm"
  configurations = [
    { type = "COMPLETE", subtype = "HOURLY", time = "2" }
  ]
}

resource "abiquo_public" "public" {
  datacenter         = "${data.abiquo_datacenter.test.id}"
  networkservicetype = "${data.abiquo_nst.test.id}"
  tag     = 2553
  mask    = 24
  address = "17.12.17.0"
  gateway = "17.12.17.1"
  name    = "test vm public"
}

resource "abiquo_external" "external" {
  enterprise         = "${data.abiquo_enterprise.test.id}"
  datacenter         = "${data.abiquo_datacenter.test.id}"
  networkservicetype = "${data.abiquo_nst.test.id}"
  tag     = 2443
  mask    = 24
  address = "172.16.6.0"
  gateway = "172.16.6.1"
  name    = "test vm external"
}

resource "abiquo_ip" "external" {
  network = "${abiquo_external.external.id}"
  ip      = "172.16.6.30"
}

resource "abiquo_ip" "public" {
  network   = "${abiquo_public.public.id}"
  ip        = "17.12.17.30"
}

resource "abiquo_virtualdatacenter" "test" {
  enterprise = "${data.abiquo_enterprise.test.id}"
  location   = "${data.abiquo_location.test.id}"
  name       = "test vm"
  type       = "KVM"
  net_mask   = "24"
  net_name   = "test vm"
  net_address = "192.168.0.0"
  net_gateway = "192.168.0.1"
  publicips  = [
    "${abiquo_ip.public.ip}"
  ]
}

resource "abiquo_firewallpolicy" "test" {
  device            = "${abiquo_virtualdatacenter.test.device}"
  virtualdatacenter = "${abiquo_virtualdatacenter.test.id}"
  # XXX workaround ABICLOUDPREMIUM-11853
  # description       = "test vm"
  name              = "test vm"
  # XXX workaround ABICLOUDPREMIUM-9668
  rules = [
    { protocol = "TCP", fromport = 22, toport = 22, sources = ["0.0.0.0/0"] }
  ]
}

resource "abiquo_private" "test" {
  virtualdatacenter = "${abiquo_virtualdatacenter.test.id}"
  mask    = 24
  address = "172.16.37.0"
  gateway = "172.16.37.1"
  name    = "test vm private"
}

resource "abiquo_ip" "private" {
  network = "${abiquo_private.test.id}"
  ip      = "172.16.37.30"
}

resource "abiquo_loadbalancer" "test" {
  device            = "${abiquo_virtualdatacenter.test.device}"
  privatenetwork    = "${abiquo_private.test.id}"
  name              = "test-vm"
  internal          = true
  algorithm         = "ROUND_ROBIN"
  routingrules      = [
    { protocolin = "HTTP" , protocolout = "HTTP" , portin = 80 , portout = 80 }
  ]
}

resource "abiquo_virtualappliance" "test" {
  virtualdatacenter = "${abiquo_virtualdatacenter.test.id}"
  name              = "test vm"
}

data "abiquo_backup" "test" {
  code     = "${abiquo_backuppolicy.test.code}"
  location = "${data.abiquo_location.test.id}"
}

data "abiquo_ip" "public" {
  pool = "${abiquo_virtualdatacenter.test.purchased}"
  ip   = "${abiquo_ip.public.ip}"
}

data "abiquo_ip" "external" {
  pool = "${abiquo_virtualdatacenter.test.externalips}"
  ip   = "${abiquo_ip.external.ip}"
}
