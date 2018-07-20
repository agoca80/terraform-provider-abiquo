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
			"abiquo_alarm":       resourceAlarm,
			"abiquo_alert":       resourceAlert,
			"abiquo_backup":      resourceBackup,
			"abiquo_computeload": resourceComputeLoad,
			"abiquo_costcode":    resourceCostcode,
			"abiquo_currency":    resourceCurrency,
			"abiquo_datacenter":  resourceDatacenter,
			"abiquo_device":      resourceDevice,
			"abiquo_dstier":      resourceDstier,
			"abiquo_enterprise":  resourceEnterprise,
			"abiquo_external":    resourceExternal,
			"abiquo_fitpolicy":   resourceFitPolicy,
			"abiquo_fw":          resourceFw,
			"abiquo_hd":          resourceHd,
			"abiquo_hp":          resourceHp,
			"abiquo_ip":          resourceIp,
			"abiquo_lb":          resourceLb,
			"abiquo_license":     resourceLicense,
			"abiquo_limit":       resourceLimit,
			"abiquo_machine":     resourceMachine,
			"abiquo_plan":        resourcePlan,
			"abiquo_pricing":     resourcePricing,
			"abiquo_private":     resourcePrivate,
			"abiquo_public":      resourcePublic,
			"abiquo_rack":        resourceRack,
			"abiquo_role":        resourceRole,
			"abiquo_scope":       resourceScope,
			"abiquo_sg":          resourceSg,
			"abiquo_storageload": resourceStorageLoad,
			"abiquo_storage":     resourceStorage,
			"abiquo_user":        resourceUser,
			"abiquo_vapp":        resourceVapp,
			"abiquo_vdc":         resourceVdc,
			"abiquo_vm":          resourceVm,
			"abiquo_vmt":         resourceVmt,
			"abiquo_vol":         resourceVolume,
		},

		DataSourcesMap: map[string]*schema.Resource{
			"abiquo_backup":     dataBackup,
			"abiquo_datacenter": dataDatacenter,
			"abiquo_devicetype": dataDeviceType,
			"abiquo_dstier":     dataDstier,
			"abiquo_enterprise": dataEnterprise,
			"abiquo_hp":         dataHp,
			"abiquo_ip":         dataIp,
			"abiquo_location":   dataLocation,
			"abiquo_machine":    dataMachine,
			"abiquo_network":    dataNetwork,
			"abiquo_nst":        dataNst,
			"abiquo_privilege":  dataPrivilege,
			"abiquo_repo":       dataRepo,
			"abiquo_role":       dataRole,
			"abiquo_scope":      dataScope,
			"abiquo_tier":       dataTier,
			"abiquo_vdc":        dataVdc,
			"abiquo_vapp":       dataVapp,
			"abiquo_template":   dataTemplate,
		},
	}
}
