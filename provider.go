package main

import (
	"sync"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

type provider struct {
	user struct {
		sync.Once
		*abiquo.User
	}

	enterprise struct {
		sync.Once
		*abiquo.Enterprise
	}
}

func (p *provider) User() *abiquo.User {
	p.user.Do(func() {
		p.user.User = abiquo.Login()
	})
	return p.user.User
}

func (p *provider) Enterprise() *abiquo.Enterprise {
	p.enterprise.Do(func() {
		p.enterprise.Enterprise = p.User().Enterprise()
	})
	return p.enterprise.Enterprise
}

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
	return new(provider), abiquo.Abiquo(d.Get("endpoint").(string), credentials)
}

// Provider factory
func Provider() *schema.Provider {
	basicAuthOptions := []string{"username", "password"}
	oAuthOptions := []string{"tokensecret", "token", "consumerkey", "consumersecret"}

	return &schema.Provider{
		ConfigureFunc: configureProvider,

		Schema: map[string]*schema.Schema{
			"endpoint": &schema.Schema{Type: schema.TypeString, Required: true, ValidateFunc: validateURL},

			"username": &schema.Schema{Type: schema.TypeString, Optional: true, ConflictsWith: oAuthOptions},
			"password": &schema.Schema{Type: schema.TypeString, Optional: true, ConflictsWith: oAuthOptions},

			"token":          &schema.Schema{Type: schema.TypeString, Optional: true, ConflictsWith: basicAuthOptions},
			"tokensecret":    &schema.Schema{Type: schema.TypeString, Optional: true, ConflictsWith: basicAuthOptions},
			"consumerkey":    &schema.Schema{Type: schema.TypeString, Optional: true, ConflictsWith: basicAuthOptions},
			"consumersecret": &schema.Schema{Type: schema.TypeString, Optional: true, ConflictsWith: basicAuthOptions},
		},

		ResourcesMap: map[string]*schema.Resource{
			"abiquo_alarm":       alarmResource,
			"abiquo_alert":       alertResource,
			"abiquo_computeload": computeLoadResource,
			"abiquo_datacenter":  datacenterResource,
			"abiquo_enterprise":  enterpriseResource,
			"abiquo_external":    externalResource,
			"abiquo_fitpolicy":   fitPolicyResource,
			"abiquo_fw":          firewallResource,
			"abiquo_hd":          hdResource,
			"abiquo_hp":          hpResource,
			"abiquo_lb":          lbResource,
			"abiquo_license":     licenseResource,
			"abiquo_limit":       limitResource,
			"abiquo_machine":     machineResource,
			"abiquo_ip":          ipResource,
			"abiquo_plan":        actionPlanResource,
			"abiquo_private":     privateResource,
			"abiquo_public":      publicResource,
			"abiquo_rack":        rackResource,
			"abiquo_role":        roleResource,
			"abiquo_scope":       scopeResource,
			"abiquo_sg":          sgResource,
			"abiquo_storageload": storageLoadResource,
			"abiquo_user":        userResource,
			"abiquo_vapp":        vappResource,
			"abiquo_vdc":         vdcResource,
			"abiquo_vm":          vmResource,
			"abiquo_vmt":         vmtResource,
			"abiquo_vol":         volResource,
		},

		DataSourcesMap: map[string]*schema.Resource{
			"abiquo_datacenter": datacenterDataSource,
			"abiquo_dstier":     dstierDataSource,
			"abiquo_enterprise": enterpriseDataSource,
			"abiquo_hp":         hpDataSource,
			"abiquo_location":   locationDataSource,
			"abiquo_machine":    machineDataSource,
			"abiquo_nst":        nstDataSource,
			"abiquo_privilege":  privilegeDataSource,
			"abiquo_repo":       repoDataSource,
			"abiquo_scope":      scopeDataSource,
			"abiquo_vdc":        vdcDataSource,
			"abiquo_template":   templateDataSource,
		},
	}
}
