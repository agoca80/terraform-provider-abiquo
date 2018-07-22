package main

var validateMedia = map[string][]string{
	// admin
	"role": []string{"/admin/roles/[0-9]+"},
	// admin/datacenters
	"backuppolicy_dc":    []string{"/admin/datacenters/[0-9]+/backuppolicies/[0-9]+"},
	"datacenter":         []string{"/admin/datacenters/[0-9]+"},
	"datastoretier_dc":   []string{"/admin/datacenters/[0-9]+/datastoretiers/[0-9]+"},
	"hardwareprofile_dc": []string{"/admin/datacenters/[0-9]+/hardwareprofiles/[0-9]+"},
	"nst":                []string{"/admin/datacenters/[0-9]+/networkservicetypes/[0-9]+"},
	"rack":               []string{"/admin/datacenters/[0-9]+/racks/[0-9]+"},
	"storagedevice":      []string{"/admin/datacenters/[0-9]+/storage/devices/[0-9]+"},
	// admin/enterprises
	"dcrepository": []string{"/admin/enterprises/[0-9]+/datacenterrepositories/[0-9]+"},
	"enterprise":   []string{"/admin/enterprises/[0-9]+"},
	"template":     []string{"/admin/enterprises/[0-9]+/datacenterrepositories/[0-9]+/virtualmachinetemplates/[0-9]+"},
	// cloud/locations
	"location":         []string{"/cloud/locations/[0-9]+"},
	"backuppolicy_vdc": []string{"/cloud/locations/[0-9]+/backuppolicies/[0-9]+"},
	"device":           []string{"/cloud/locations/[0-9]+/devices/[0-9]+"},
	"firewall":         []string{"/cloud/locations/[0-9]+/devices/[0-9]+/firewalls/[0-9]+"},
	"loadbalancer":     []string{"/cloud/locations/[0-9]+/devices/[0-9]+/loadbalancers/[0-9]+"},
	// cloud/virtualdatacenters
	"privatenetwork":    []string{"/cloud/virtualdatacenters/[0-9]+/privatenetworks/[0-9]+"},
	"tier_vdc":          []string{"/cloud/virtualdatacenters/[0-9]+/tiers/[0-9]+"},
	"virtualappliance":  []string{"/cloud/virtualdatacenters/[0-9]+/virtualappliances/[0-9]+"},
	"virtualdatacenter": []string{"/cloud/virtualdatacenters/[0-9]+"},
	// Virtualmachine
	"alarm":          []string{"/cloud/virtualdatacenters/[0-9]+/virtualappliances/[0-9]+/virtualmachines/[0-9]+/metrics/.+/alarms/.+"},
	"virtualmachine": []string{"/cloud/virtualdatacenters/[0-9]+/virtualappliances/[0-9]+/virtualmachines/[0-9]+"},
	"virtualmachine_ip": []string{
		"/admin/enterprises/[0-9]+/limits/[0-9]+/externalnetworks/[0-9]+/ips/[0-9]+",
		"/cloud/virtualdatacenters/[0-9]+/privatenetworks/[0-9]+/ips/[0-9]+",
		"/cloud/virtualdatacenters/[0-9]+/publicips/purchased/[0-9]+",
	},
	"templates": []string{
		"/admin/enterprises/[0-9]+/datacenterrepositories/[0-9]+/virtualmachinetemplates",
		"/cloud/virtualdatacenters/[0-9]+/action/templates",
	},
	// config
	"currency":   []string{"/config/currencies/[0-9]+"},
	"devicetype": []string{"/config/devicetypes/.+"},
}
