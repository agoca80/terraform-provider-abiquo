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
			"endpoint": &schema.Schema{
				DefaultFunc:  schema.EnvDefaultFunc("ABQ_ENDPOINT", ""),
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validateURL,
			},
			"username": &schema.Schema{
				ConflictsWith: oAuthOptions,
				DefaultFunc:   schema.EnvDefaultFunc("ABQ_USERNAME", ""),
				Optional:      true,
				Type:          schema.TypeString,
			},
			"password": &schema.Schema{
				ConflictsWith: oAuthOptions,
				DefaultFunc:   schema.EnvDefaultFunc("ABQ_PASSWORD", ""),
				Optional:      true,
				Type:          schema.TypeString,
			},
			"token":          &schema.Schema{Type: schema.TypeString, Optional: true, ConflictsWith: basicAuthOptions},
			"tokensecret":    &schema.Schema{Type: schema.TypeString, Optional: true, ConflictsWith: basicAuthOptions},
			"consumerkey":    &schema.Schema{Type: schema.TypeString, Optional: true, ConflictsWith: basicAuthOptions},
			"consumersecret": &schema.Schema{Type: schema.TypeString, Optional: true, ConflictsWith: basicAuthOptions},
		},

		ResourcesMap: map[string]*schema.Resource{
			"abiquo_alarm":       alarmResource,
			"abiquo_alert":       alertResource,
			"abiquo_backup":      backupResource,
			"abiquo_computeload": computeLoadResource,
			"abiquo_datacenter":  datacenterResource,
			"abiquo_device":      deviceResource,
			"abiquo_dstier":      dstierResource,
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
			"abiquo_backup": &schema.Resource{
				Schema: backupDataSchema,
				Read:   backupDataRead,
			},
			"abiquo_datacenter": &schema.Resource{
				Schema: datacenterDataSchema,
				Read:   datacenterDataRead,
			},
			"abiquo_devicetype": &schema.Resource{
				Schema: deviceTypeDataSchema,
				Read:   deviceTypeDataRead,
			},
			"abiquo_dstier": &schema.Resource{
				Schema: dstierDataSchema,
				Read:   dstierDataRead,
			},
			"abiquo_enterprise": &schema.Resource{
				Schema: enterpriseDataSchema,
				Read:   enterpriseDataRead,
			},
			"abiquo_hp": &schema.Resource{
				Schema: hpDataSchema,
				Read:   hpDataRead,
			},
			"abiquo_location": &schema.Resource{
				Schema: locationDataSchema,
				Read:   locationRead,
			},
			"abiquo_machine": &schema.Resource{
				Schema: machineDataSchema,
				Read:   machineDataRead,
			},
			"abiquo_nst": &schema.Resource{
				Schema: nstDataSchema,
				Read:   nstDataRead,
			},
			"abiquo_privilege": &schema.Resource{
				Schema: privilegeDataSchema,
				Read:   privilegeRead,
			},
			"abiquo_repo": &schema.Resource{
				Schema: repoDataSchema,
				Read:   dataRepoRead,
			},
			"abiquo_role": &schema.Resource{
				Schema: roleDataSchema,
				Read:   roleDataRead,
			},
			"abiquo_scope": &schema.Resource{
				Schema: scopeDataSchema,
				Read:   scopeDataRead,
			},
			"abiquo_vdc": &schema.Resource{
				Schema: vdcDataSchema,
				Read:   dataVDCRead,
			},
			"abiquo_vapp": &schema.Resource{
				Schema: vappDataSchema,
				Read:   vappDataRead,
			},
			"abiquo_template": &schema.Resource{
				Schema: templateDataSchema,
				Read:   templateRead,
			},
		},
	}
}
