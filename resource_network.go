package main

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
)

func networkNew(d *resourceData) *abiquo.Network {
	return &abiquo.Network{
		Address: d.string("address"),
		Mask:    d.integer("mask"),
		Gateway: d.string("gateway"),
		Name:    d.string("name"),
		DNS1:    d.string("dns1"),
		DNS2:    d.string("dns2"),
		Suffix:  d.string("suffix"),
	}
}

func networkRead(d *resourceData, resource core.Resource) {
	network := resource.(*abiquo.Network)
	d.Set("tag", network.Tag)
	d.Set("address", network.Address)
	d.Set("mask", network.Mask)
	d.Set("gateway", network.Gateway)
	d.Set("dns1", network.DNS1)
	d.Set("dns2", network.DNS2)
	d.Set("suffix", network.Suffix)
	return
}
