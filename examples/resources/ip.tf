resource "abiquo_private" "private" {
  virtualdatacenter = "${abiquo_vdc.vdc.id}"

  # XXX workaround ABICLOUDPREMIUM-9660
  lifecycle = { ignore_changes = [ "dns1", "dns2" ] }

  mask    = 24
  address = "172.16.0.0"
  gateway = "172.16.0.1"
  name    = "terraform private"
  dns1    = "8.8.8.8"
  dns2    = "4.4.4.4"
}

resource "abiquo_ip" "private" {
  ip      = "172.16.0.3"
  type    = "privateip"
  network = "${abiquo_private.private.id}"
}

# resource "abiquo_ip" "external" {
#   ip      = "172.16.4.10"
#   type    = "externalip"
#   network = "${abiquo_external.example.id}"
# }
# 
# resource "abiquo_ip" "public" {
#   ip      = "172.16.178.3"
#   type    = "publicip"
#   network = "${abiquo_public.public.id}"
# }

resource "abiquo_vapp" "vapp" {
  virtualdatacenter = "${abiquo_vdc.vdc.id}"  
  name = "terraform vapp"
}

# resource "abiquo_vm" "ips" {
#   label     = "terraform ips"
#   ips = [
#     "${abiquo_ip.private.id}",
#     "${abiquo_ip.external.id}",
#     "${abiquo_ip.public.id}"
#   ]
#   virtualappliance       = "${abiquo_vapp.vapp.id}"
#   virtualmachinetemplate = "${var.template}"
# }