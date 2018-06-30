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
			"abiquo_alarm": &schema.Resource{
				Schema: alarmSchema,
				Delete: resourceDelete,
				Exists: resourceExists("alarm"),
				Update: resourceUpdate(alarmNew, nil, "alarm"),
				Create: resourceCreate(alarmNew, nil, alarmRead, alarmEndpoint),
				Read:   resourceRead(alarmNew, alarmRead, "alarm"),
			},
			"abiquo_alert": &schema.Resource{
				Schema: alertSchema,
				Delete: resourceDelete,
				Exists: resourceExists("alert"),
				Update: resourceUpdate(alertNew, nil, "alert"),
				Create: resourceCreate(alertNew, nil, alertRead, alertEndpoint),
				Read:   resourceRead(alertNew, alertRead, "alert"),
			},
			"abiquo_backup": &schema.Resource{
				Schema: backupSchema,
				Read:   resourceRead(backupDTO, backupRead, "backuppolicy"),
				Update: resourceUpdate(backupDTO, nil, "backuppolicy"),
				Exists: resourceExists("backuppolicy"),
				Delete: resourceDelete,
				Create: resourceCreate(backupDTO, nil, backupRead, backupEndpoint),
			},
			"abiquo_computeload": &schema.Resource{
				Schema: computeLoadSchema,
				Delete: resourceDelete,
				Exists: resourceExists("machineloadrule"),
				Create: resourceCreate(computeLoadDTO, nil, computeLoadRead, computeLoadEndpoint),
				Update: resourceUpdate(computeLoadDTO, nil, "machineloadrule"),
				Read:   resourceRead(computeLoadDTO, computeLoadRead, "machineloadrule"),
			},
			"abiquo_costcode": &schema.Resource{
				Schema: costCodeSchema,
				Read:   resourceRead(costCodeNew, costCodeRead, "costcode"),
				Update: resourceUpdate(costCodeNew, nil, "costcode"),
				Exists: resourceExists("costcode"),
				Delete: resourceDelete,
				Create: resourceCreate(costCodeNew, nil, costCodeRead, costCodeEndpoint),
			},
			"abiquo_currency": &schema.Resource{
				Schema: currencySchema,
				Read:   resourceRead(currencyNew, currencyRead, "currency"),
				Update: resourceUpdate(currencyNew, nil, "currency"),
				Exists: resourceExists("currency"),
				Delete: resourceDelete,
				Create: resourceCreate(currencyNew, nil, currencyRead, currencyEndpoint),
			},
			"abiquo_datacenter": &schema.Resource{
				Schema: datacenterSchema,
				Delete: resourceDelete,
				Exists: resourceExists("datacenter"),
				Create: resourceCreate(datacenterNew, nil, datacenterRead, datacenterEndpoint),
				Update: resourceUpdate(datacenterNew, nil, "datacenter"),
				Read:   resourceRead(datacenterNew, datacenterRead, "datacenter"),
			},
			"abiquo_device": &schema.Resource{
				Schema: deviceSchema,
				Delete: resourceDelete,
				Exists: resourceExists("device"),
				Create: resourceCreate(deviceDTO, nil, deviceRead, deviceEndpoint),
				Update: resourceUpdate(deviceDTO, nil, "device"),
				Read:   resourceRead(deviceDTO, deviceRead, "device"),
			},
			"abiquo_dstier": &schema.Resource{
				Schema: dstierSchema,
				Delete: resourceDelete,
				Exists: resourceExists("datastoretier"),
				Create: resourceCreate(dstierDTO, nil, dstierRead, dstierEndpoint),
				Update: resourceUpdate(dstierDTO, nil, "datastoretier"),
				Read:   resourceRead(dstierDTO, dstierRead, "datastoretier"),
			},
			"abiquo_enterprise": &schema.Resource{
				Schema: enterpriseSchema,
				Delete: resourceDelete,
				Read:   resourceRead(enterpriseDTO, enterpriseRead, "enterprise"),
				Create: resourceCreate(enterpriseDTO, enterpriseCreate, enterpriseRead, enterpriseEndpoint),
				Exists: resourceExists("enterprise"),
				Update: resourceUpdate(enterpriseDTO, enterpriseUpdate, "enterprise"),
			},
			"abiquo_external": &schema.Resource{
				Schema: externalSchema,
				Delete: resourceDelete,
				Exists: resourceExists("vlan"),
				Update: resourceUpdate(externalNew, nil, "vlan"),
				Create: resourceCreate(externalNew, nil, externalRead, externalEndpoint),
				Read:   resourceRead(externalNew, externalRead, "vlan"),
			},
			"abiquo_fitpolicy": &schema.Resource{
				Schema: fitPolicySchema,
				Delete: resourceDelete,
				Exists: resourceExists("fitpolicyrule"),
				Create: resourceCreate(fitPolicyDTO, nil, fitPolicyRead, fitPolicyEndpoint),
				Read:   resourceRead(fitPolicyDTO, fitPolicyRead, "fitpolicyrule"),
			},
			"abiquo_fw": &schema.Resource{
				Schema: firewallSchema,
				Delete: resourceDelete,
				Exists: resourceExists("firewallpolicy"),
				Update: resourceUpdate(fwNew, fwUpdate, "firewallpolicy"),
				Create: resourceCreate(fwNew, fwCreate, fwRead, fwEndpoint),
				Read:   resourceRead(fwNew, fwRead, "firewallpolicy"),
			},
			"abiquo_hd": &schema.Resource{
				Schema: hdSchema,
				Update: schema.Noop,
				Delete: schema.Noop,
				Create: resourceCreate(hdNew, nil, hdRead, hdEndpoint),
				Exists: resourceExists("harddisk"),
				Read:   resourceRead(hdNew, hdRead, "harddisk"),
			},
			"abiquo_hp": &schema.Resource{
				Schema: hpSchema,
				Delete: resourceDelete,
				Exists: resourceExists("hardwareprofile"),
				Create: resourceCreate(hpNew, nil, hpRead, hpEndpoint),
				Update: resourceUpdate(hpNew, nil, "hardwareprofile"),
				Read:   resourceRead(hpNew, hpRead, "hardwareprofile"),
			},
			"abiquo_ip": &schema.Resource{
				Schema: ipSchema,
				Delete: resourceDelete,
				Exists: ipExists,
				Create: ipCreate,
				Read:   ipRead,
			},
			"abiquo_lb": &schema.Resource{
				Schema: lbSchema,
				Delete: resourceDelete,
				Exists: resourceExists("loadbalancer"),
				Create: resourceCreate(lbNew, nil, lbRead, lbEndpoint),
				Update: resourceUpdate(lbNew, nil, "loadbalancer"),
				Read:   resourceRead(lbNew, lbRead, "loadbalancer"),
			},
			"abiquo_license": &schema.Resource{
				Schema: licenseSchema,
				Delete: resourceDelete,
				Exists: resourceExists("license"),
				Create: resourceCreate(licenseNew, nil, licenseRead, licenseEndpoint),
				Read:   resourceRead(licenseNew, licenseRead, "license"),
			},
			"abiquo_limit": &schema.Resource{
				Schema: limitSchema,
				Exists: resourceExists("limit"),
				Read:   resourceRead(limitNew, limitRead, "limit"),
				Update: resourceUpdate(limitNew, nil, "limit"),
				Create: resourceCreate(limitNew, nil, limitRead, limitEndpoint),
				Delete: resourceDelete,
			},
			"abiquo_machine": &schema.Resource{
				Schema: machineSchema,
				Delete: resourceDelete,
				Exists: resourceExists("machine"),
				Create: machineCreate,
				Update: machineUpdate,
				Read:   machineRead,
			},
			"abiquo_plan": &schema.Resource{
				Schema: actionPlanSchema,
				Delete: resourceDelete,
				Exists: resourceExists("virtualmachineactionplan"),
				Update: resourceUpdate(actionPlanNew, nil, "virtualmachineactionplan"),
				Create: resourceCreate(actionPlanNew, actionPlanCreate, actionPlanRead, actionPlanEndpoint),
				Read:   resourceRead(actionPlanNew, actionPlanRead, "virtualmachineactionplan"),
			},
			"abiquo_pricing": &schema.Resource{
				Schema: pricingSchema,
				Delete: resourceDelete,
				Exists: resourceExists("pricingtemplate"),
				Update: resourceUpdate(pricingNew, nil, "pricingtemplate"),
				Create: resourceCreate(pricingNew, nil, pricingRead, pricingEndpoint),
				Read:   resourceRead(pricingNew, pricingRead, "pricingtemplate"),
			},
			"abiquo_private": &schema.Resource{
				Schema: privateSchema,
				Delete: resourceDelete,
				Update: resourceUpdate(privateNew, nil, "vlan"),
				Create: resourceCreate(privateNew, nil, privateRead, privateEndpoint),
				Read:   resourceRead(privateNew, privateRead, "vlan"),
			},
			"abiquo_public": &schema.Resource{
				Schema: publicSchema,
				Delete: resourceDelete,
				Update: resourceUpdate(publicNew, nil, "vlan"),
				Create: resourceCreate(publicNew, nil, publicRead, publicEndpoint),
				Read:   resourceRead(publicNew, publicRead, "vlan"),
			},
			"abiquo_rack": &schema.Resource{
				Schema: rackSchema,
				Delete: resourceDelete,
				Exists: resourceExists("rack"),
				Create: resourceCreate(rackNew, nil, rackRead, rackEndpoint),
				Update: resourceUpdate(rackNew, nil, "rack"),
				Read:   resourceRead(rackNew, rackRead, "rack"),
			},
			"abiquo_role": &schema.Resource{
				Schema: roleSchema,
				Delete: resourceDelete,
				Read:   resourceRead(roleNew, roleRead, "role"),
				Create: resourceCreate(roleNew, rolePrivileges, roleRead, roleEndpoint),
				Exists: resourceExists("role"),
				Update: resourceUpdate(roleNew, rolePrivileges, "role"),
			},
			"abiquo_scope": &schema.Resource{
				Schema: scopeSchema,
				Delete: resourceDelete,
				Read:   resourceRead(scopeNew, scopeRead, "scope"),
				Create: resourceCreate(scopeNew, nil, scopeRead, scopeEndpoint),
				Exists: resourceExists("scope"),
				Update: resourceUpdate(scopeNew, nil, "scope"),
			},
			"abiquo_sg": &schema.Resource{
				Schema: sgSchema,
				Delete: sgDelete,
				Update: resourceUpdate(sgNew, nil, "scalinggroup"),
				Create: resourceCreate(sgNew, nil, sgRead, sgEndpoint),
				Read:   resourceRead(sgNew, sgRead, "scalinggroup"),
			},
			"abiquo_storageload": &schema.Resource{
				Schema: storageLoadSchema,
				Delete: resourceDelete,
				Exists: resourceExists("datastoreloadrule"),
				Create: resourceCreate(storageLoadDTO, nil, storageLoadRead, storageLoadEndpoint),
				Update: resourceUpdate(storageLoadDTO, nil, "datastoreloadrule"),
				Read:   resourceRead(storageLoadDTO, storageLoadRead, "datastoreloadrule"),
			},
			"abiquo_storage": &schema.Resource{
				Schema: storageDeviceSchema,
				Delete: resourceDelete,
				Exists: resourceExists("virtualappliance"),
				Create: resourceCreate(storageDeviceNew, nil, storageDeviceRead, storageDeviceEndpoint),
				Update: resourceUpdate(storageDeviceNew, nil, "storagedevice"),
				Read:   resourceRead(storageDeviceNew, storageDeviceRead, "storagedevice"),
			},
			"abiquo_user": &schema.Resource{
				Schema: userSchema,
				Read:   resourceRead(userNew, userRead, "user"),
				Create: resourceCreate(userNew, nil, userRead, userEndpoint),
				Update: resourceUpdate(userNew, nil, "user"),
				Exists: resourceExists("user"),
				Delete: resourceDelete,
			},
			"abiquo_vapp": &schema.Resource{
				Schema: vappSchema,
				Delete: resourceDelete,
				Exists: resourceExists("virtualappliance"),
				Create: resourceCreate(vappNew, nil, vappRead, vappEndpoint),
				Update: resourceUpdate(vappNew, nil, "virtualappliance"),
				Read:   resourceRead(vappNew, vappRead, "virtualappliance"),
			},
			"abiquo_vdc": &schema.Resource{
				Schema: vdcSchema,
				Delete: resourceDelete,
				Create: resourceCreate(vdcNew, vdcCreate, vdcRead, vdcEndpoint),
				Exists: resourceExists("virtualdatacenter"),
				Update: resourceUpdate(vdcNew, vdcUpdate, "virtualdatacenter"),
				Read:   resourceRead(vdcNew, vdcRead, "virtualdatacenter"),
			},
			"abiquo_vm": &schema.Resource{
				Schema: vmSchema,
				Read:   resourceRead(vmNew, vmRead, "virtualmachine"),
				Exists: resourceExists("virtualmachine"),
				Delete: vmDelete,
				Create: resourceCreate(vmNew, vmCreate, vmRead, vmEndpoint),
			},
			"abiquo_vmt": &schema.Resource{
				Schema: vmtSchema,
				Create: vmtCreate,
				Delete: resourceDelete,
				Update: vmtUpdate,
				Read:   resourceRead(vmtNew, vmtRead, "virtualmachinetemplate"),
				Exists: resourceExists("virtualmachinetemplate"),
			},
			"abiquo_vol": &schema.Resource{
				Schema: volumeSchema,
				Delete: resourceDelete,
				Update: resourceUpdate(volNew, nil, "volume"),
				Create: resourceCreate(volNew, nil, volRead, volEndpoint),
				Read:   resourceRead(volNew, volRead, "volume"),
			},
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
			"abiquo_ip": &schema.Resource{
				Schema: ipDataSchema,
				Read:   ipDataRead,
			},
			"abiquo_location": &schema.Resource{
				Schema: locationDataSchema,
				Read:   locationRead,
			},
			"abiquo_machine": &schema.Resource{
				Schema: machineDataSchema,
				Read:   machineDataRead,
			},
			"abiquo_network": &schema.Resource{
				Schema: networkDataSchema,
				Read:   networkDataRead,
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
			"abiquo_tier": &schema.Resource{
				Schema: tierDataSchema,
				Read:   tierDataRead,
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
