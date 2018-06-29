data "abiquo_enterprise" "test" { name = "Abiquo" }
data "abiquo_location"   "test" { name = "datacenter 1" }
data "abiquo_datacenter" "test" { name = "datacenter 1" }
data "abiquo_template"   "test" { name = "tests" }
data "abiquo_nst"        "test"        {
  datacenter = "${data.abiquo_datacenter.test.id}"
  name       = "Service Network"
}

resource "abiquo_backup" "test" {
  datacenter     = "${data.abiquo_datacenter.test.id}"
  code           = "testVM"
  name           = "testVM"
  description    = "testVM"
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
  name    = "testVM-public"
}

resource "abiquo_external" "external" {
  enterprise         = "${data.abiquo_enterprise.test.id}"
  datacenter         = "${data.abiquo_datacenter.test.id}"
  networkservicetype = "${data.abiquo_nst.test.id}"

  tag     = 2443
  mask    = 24
  address = "172.16.6.0"
  gateway = "172.16.6.1"
  name    = "testVM-external"
}

resource "abiquo_ip" "external" {
  network = "${abiquo_external.external.id}"
  ip      = "172.16.6.30"
}

resource "abiquo_ip" "public" {
  network   = "${abiquo_public.public.id}"
  ip        = "17.12.17.30"
}

resource "abiquo_vdc" "test" {
  enterprise = "${data.abiquo_enterprise.test.id}"
  location   = "${data.abiquo_location.test.id}"
  name       = "testVM"
  type       = "KVM"
  publicips  = [
    "${abiquo_ip.public.ip}"
  ]
}

resource "abiquo_fw" "test" {
  virtualdatacenter = "${abiquo_vdc.test.id}"

  description = "testVM"
  name        = "testVM"

  # XXX workaround ABICLOUDPREMIUM-9668
  rules = [
    { protocol = "TCP", fromport = 22, toport = 22, sources = ["0.0.0.0/0"] }
  ]
}

resource "abiquo_private" "test" {
  virtualdatacenter = "${abiquo_vdc.test.id}"
  mask    = 24
  address = "172.16.37.0"
  gateway = "172.16.37.1"
  name    = "testVM-private"
}

resource "abiquo_ip" "private" {
  network = "${abiquo_private.test.id}"
  ip      = "172.16.37.30"
}

resource "abiquo_lb" "test" {
  virtualdatacenter = "${abiquo_vdc.test.id}"
  privatenetwork    = "${abiquo_private.test.id}"

  name         = "testVM"
  algorithm    = "ROUND_ROBIN"
  routingrules = [
    { protocolin = "HTTP" , protocolout = "HTTP" , portin = 80 , portout = 80 }
  ]
}

resource "abiquo_vapp" "test" {
  virtualdatacenter = "${abiquo_vdc.test.id}"
  name              = "testVM"
}

data "abiquo_backup" "test" {
  code     = "${abiquo_backup.test.code}"
  location = "${data.abiquo_location.test.id}"
}

data "abiquo_ip" "public" {
  pool = "${abiquo_vdc.test.purchased}"
  ip   = "${abiquo_ip.public.ip}"
}

data "abiquo_ip" "external" {
  pool = "${abiquo_vdc.test.externalips}"
  ip   = "${abiquo_ip.external.ip}"
}

resource "abiquo_vm" "test" {
  deploy                 = false
  backups                = [ "${data.abiquo_backup.test.id}" ]
  cpu                    = 1
  ram                    = 64
  label                  = "testVM"
  virtualappliance       = "${abiquo_vapp.test.id}"
  virtualmachinetemplate = "${data.abiquo_template.test.id}"

  lbs = [ "${abiquo_lb.test.id}" ]
  fws = [ "${abiquo_fw.test.id}" ]

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
