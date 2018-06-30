data     "abiquo_devicetype" "test" { name = "LOGICAL" }
data     "abiquo_datacenter" "test" { name = "datacenter 1" }

resource "abiquo_device"     "test" {
  default    = false
  devicetype = "${data.abiquo_devicetype.test.id}"
  endpoint   = "https://logical:35353/api"
  name       = "test device"
  username   = "username"
  password   = "password"
  datacenter = "${data.abiquo_datacenter.test.id}"
}
