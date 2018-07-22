package main

import (
	"sync"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

type provider struct {
	err        error
	init       sync.Once
	user       *abiquo.User
	enterprise *abiquo.Enterprise
}

var abq provider

func (p *provider) User() *abiquo.User             { return p.user }
func (p *provider) Enterprise() *abiquo.Enterprise { return p.enterprise }

func configureProvider(d *schema.ResourceData) (meta interface{}, err error) {
	var credentials interface{}
	if _, ok := d.GetOk("username"); ok {
		credentials = core.Basic{
			Username: d.Get("username").(string),
			Password: d.Get("password").(string),
		}
	} else {
		credentials = core.Oauth{
			APIKey:      d.Get("consumerkey").(string),
			APISecret:   d.Get("consumersecret").(string),
			Token:       d.Get("token").(string),
			TokenSecret: d.Get("tokensecret").(string),
		}
	}

	abq.init.Do(func() {
		endpoint := d.Get("endpoint").(string)
		if abq.err = core.Init(endpoint, credentials); abq.err != nil {
			return
		}
		abq.user = abiquo.Login()
		resource := abq.user.Rel("enterprise").Walk()
		abq.enterprise = resource.(*abiquo.Enterprise)
	})
	return &abq, abq.err
}

// Provider factory
func Provider() *schema.Provider {
	basicAuthOptions := []string{"username", "password"}
	oAuthOptions := []string{"tokensecret", "token", "consumerkey", "consumersecret"}

	resourceMap := make(map[string]*schema.Resource)
	for _, d := range resources {
		resourceMap[d.Name()] = d.resource()
	}

	return &schema.Provider{
		ConfigureFunc: configureProvider,

		Schema: map[string]*schema.Schema{
			"endpoint":       attribute(required, href, variable("ABQ_ENDPOINT")),
			"username":       attribute(optional, text, variable("ABQ_USERNAME"), conflicts(oAuthOptions)),
			"password":       attribute(optional, text, variable("ABQ_PASSWORD"), conflicts(oAuthOptions)),
			"token":          attribute(optional, text, conflicts(basicAuthOptions)),
			"tokensecret":    attribute(optional, text, conflicts(basicAuthOptions)),
			"consumerkey":    attribute(optional, text, conflicts(basicAuthOptions)),
			"consumersecret": attribute(optional, text, conflicts(basicAuthOptions)),
		},

		ResourcesMap: resourceMap,

		DataSourcesMap: map[string]*schema.Resource{
			"abiquo_backup":            dataSource(backupDataSchema, backupFind),
			"abiquo_datacenter":        dataSource(datacenterDataSchema, datacenterFind),
			"abiquo_devicetype":        dataSource(deviceTypeDataSchema, deviceTypeFind),
			"abiquo_datastoretier":     dataSource(dstierDataSchema, dstierFind),
			"abiquo_enterprise":        dataSource(enterpriseDataSchema, enterpriseFind),
			"abiquo_hardwareprofile":   dataSource(hpDataSchema, hpFind),
			"abiquo_ip":                dataSource(ipDataSchema, ipFind),
			"abiquo_location":          dataSource(locationDataSchema, locationFind),
			"abiquo_machine":           dataSource(machineDataSchema, machineFind),
			"abiquo_network":           dataSource(networkDataSchema, networkFind),
			"abiquo_nst":               dataSource(nstDataSchema, nstFind),
			"abiquo_privilege":         dataSource(privilegeDataSchema, privilegeFind),
			"abiquo_repo":              dataSource(repoDataSchema, repoFind),
			"abiquo_role":              dataSource(roleDataSchema, roleFind),
			"abiquo_scope":             dataSource(scopeDataSchema, scopeFind),
			"abiquo_tier":              dataSource(tierDataSchema, tierFind),
			"abiquo_virtualappliance":  dataSource(vappDataSchema, vappFind),
			"abiquo_virtualdatacenter": dataSource(vdcDataSchema, virtualdatacenterFind),
			"abiquo_template":          dataSource(templateDataSchema, templateFind),
		},
	}
}

var resources = []*description{
	alarm, alert, backuppolicy, machineloadrule, costcode, currency, datacenter,
	device, datastoretier, enterprise, external, fitpolicyrule, firewallpolicy,
	harddisk, hardwareprofile, ipAddress, loadbalancer, license, limit, machine,
	virtualmachineactionplan, pricingtemplate, private, public, rack, role,
	scope, scalinggroup, datastoreloadrule, storagedevice, user,
	virtualappliance, virtualdatacenter, virtualmachine, virtualmachinetemplate,
	volume,
}
