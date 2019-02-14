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
	var credentials core.Authenticator
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
		abq.user, err = abiquo.Login()
		if err != nil {
			return
		}

		resource, err := abq.user.Rel("enterprise").Walk()
		if err != nil {
			return
		}
		abq.enterprise = resource.(*abiquo.Enterprise)
	})
	return &abq, abq.err
}

// Provider factory
func Provider() *schema.Provider {
	basicAuthOptions := []string{"username", "password"}
	oAuthOptions := []string{"tokensecret", "token", "consumerkey", "consumersecret"}

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

		ResourcesMap: map[string]*schema.Resource{
			"abiquo_actionplan":             resourceDefinition(actionplan),
			"abiquo_alarm":                  resourceDefinition(alarm),
			"abiquo_alert":                  resourceDefinition(alert),
			"abiquo_backuppolicy":           resourceDefinition(backuppolicy),
			"abiquo_costcode":               resourceDefinition(costcode),
			"abiquo_currency":               resourceDefinition(currency),
			"abiquo_datacenter":             resourceDefinition(datacenter),
			"abiquo_datastoreloadrule":      resourceDefinition(datastoreloadrule),
			"abiquo_device":                 resourceDefinition(device),
			"abiquo_datastoretier":          resourceDefinition(datastoretier),
			"abiquo_enterprise":             resourceDefinition(enterprise),
			"abiquo_external":               resourceDefinition(external),
			"abiquo_fitpolicyrule":          resourceDefinition(fitpolicyrule),
			"abiquo_firewallpolicy":         resourceDefinition(firewallpolicy),
			"abiquo_harddisk":               resourceDefinition(harddisk),
			"abiquo_hardwareprofile":        resourceDefinition(hardwareprofile),
			"abiquo_ip":                     resourceDefinition(ipAddress),
			"abiquo_loadbalancer":           resourceDefinition(loadbalancer),
			"abiquo_license":                resourceDefinition(license),
			"abiquo_limit":                  resourceDefinition(limit),
			"abiquo_machine":                resourceDefinition(machine),
			"abiquo_machineloadrule":        resourceDefinition(machineloadrule),
			"abiquo_pricingtemplate":        resourceDefinition(pricingtemplate),
			"abiquo_private":                resourceDefinition(private),
			"abiquo_public":                 resourceDefinition(public),
			"abiquo_rack":                   resourceDefinition(rack),
			"abiquo_role":                   resourceDefinition(role),
			"abiquo_scope":                  resourceDefinition(scope),
			"abiquo_scalinggroup":           resourceDefinition(scalinggroup),
			"abiquo_storagedevice":          resourceDefinition(storagedevice),
			"abiquo_user":                   resourceDefinition(user),
			"abiquo_virtualappliance":       resourceDefinition(virtualappliance),
			"abiquo_virtualdatacenter":      resourceDefinition(virtualdatacenter),
			"abiquo_virtualmachine":         resourceDefinition(virtualmachine),
			"abiquo_virtualmachinetemplate": resourceDefinition(virtualmachinetemplate),
			"abiquo_volume":                 resourceDefinition(volume),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"abiquo_backup":            data(backupDataSchema, backupFind),
			"abiquo_datacenter":        data(datacenterDataSchema, datacenterFind),
			"abiquo_devicetype":        data(deviceTypeDataSchema, deviceTypeFind),
			"abiquo_datastoretier":     data(dstierDataSchema, dstierFind),
			"abiquo_enterprise":        data(enterpriseDataSchema, enterpriseFind),
			"abiquo_hardwareprofile":   data(hpDataSchema, hpFind),
			"abiquo_ip":                data(ipDataSchema, ipFind),
			"abiquo_location":          data(locationDataSchema, locationFind),
			"abiquo_machine":           data(machineDataSchema, machineFind),
			"abiquo_network":           data(networkDataSchema, networkFind),
			"abiquo_nst":               data(nstDataSchema, nstFind),
			"abiquo_privilege":         data(privilegeDataSchema, privilegeFind),
			"abiquo_repo":              data(repoDataSchema, repoFind),
			"abiquo_role":              data(roleDataSchema, roleFind),
			"abiquo_scope":             data(scopeDataSchema, scopeFind),
			"abiquo_tier":              data(tierDataSchema, tierFind),
			"abiquo_virtualappliance":  data(vappDataSchema, vappFind),
			"abiquo_virtualdatacenter": data(vdcDataSchema, virtualdatacenterFind),
			"abiquo_template":          data(templateDataSchema, templateFind),
		},
	}
}
