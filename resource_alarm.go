package main

import (
	"fmt"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/schema"
)

var alarmSchema = map[string]*schema.Schema{
	"target":      attribute(required, href, forceNew),
	"formula":     attribute(required, text),
	"name":        attribute(required, text),
	"metric":      attribute(required, text, forceNew),
	"evaluations": attribute(required, natural),
	"period":      attribute(required, natural),
	"statistic":   attribute(required, text),
	"threshold":   attribute(required, integer),
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
		EvaluationPeriods: d.integer("evaluations"),
		Name:              d.string("name"),
		Formula:           d.string("formula"),
		Period:            d.integer("period"),
		Statistic:         d.string("statistic"),
		Threshold:         d.integer("threshold"),
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
