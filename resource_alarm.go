package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var alarmResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"target":      Required().Renew().Link(),
		"formula":     Required().String(),
		"name":        Required().String(),
		"metric":      Required().Renew().String(),
		"evaluations": Required().Number(),
		"period":      Required().Number(),
		"statistic":   Required().String(),
		"threshold":   Required().Number(),
	},
	Delete: resourceDelete,
	Exists: resourceExists("alarm"),
	Update: resourceUpdate(alarmNew, nil, "alarm"),
	Create: resourceCreate(alarmNew, nil, alarmRead, alarmEndpoint),
	Read:   resourceRead(alarmNew, alarmRead, "alarm"),
}

func alarmEndpoint(d *resourceData) *core.Link {
	target := d.string("target")
	metric := d.string("metric")
	alarms := fmt.Sprintf("%v/metrics/%v/alarms", target, metric)
	return core.NewLinkType(alarms, "alarm")
}

func alarmNew(d *resourceData) core.Resource {
	target := d.string("target")
	metric := d.string("metric")
	href := fmt.Sprintf("%v/metrics/%v", target, metric)
	return &abiquo.Alarm{
		EvaluationPeriods: d.int("evaluations"),
		Name:              d.string("name"),
		Formula:           d.string("formula"),
		Period:            d.int("period"),
		Statistic:         d.string("statistic"),
		Threshold:         d.int("threshold"),
		DTO: core.NewDTO(
			core.NewLinkType(href, "metric").SetRel("metric"),
		),
	}
}

func alarmRead(d *resourceData, resource core.Resource) (err error) {
	alarm := resource.(*abiquo.Alarm)
	d.Set("name", alarm.Name)
	d.Set("evaluations", alarm.EvaluationPeriods)
	d.Set("formula", alarm.Formula)
	d.Set("period", alarm.Period)
	d.Set("statistic", alarm.Statistic)
	d.Set("threshold", alarm.Threshold)
	return
}
