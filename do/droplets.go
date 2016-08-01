/*
Copyright 2016 The Doctl Authors All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package do

import (
	"github.com/digitalocean/godo"
	"github.com/digitalocean/godo/util"
)

// DropletIPTable is a table of interface IPS.
type DropletIPTable map[InterfaceType]string

// InterfaceType is a an interface type.
type InterfaceType string

const (
	// InterfacePublic is a public interface.
	InterfacePublic InterfaceType = "public"
	// InterfacePrivate is a private interface.
	InterfacePrivate InterfaceType = "private"
)

// Droplet is a wrapper for godo.Droplet
type Droplet struct {
	*godo.Droplet
}

// Droplets is a slice of Droplet.
type Droplets []Droplet

// Kernel is a wrapper for godo.Kernel
type Kernel struct {
	*godo.Kernel
}

// Kernels is a slice of Kernel.
type Kernels []Kernel

// DropletsService is an interface for interacting with DigitalOcean's droplet api.
type DropletsService interface {
	List() (Droplets, error)
	ListByTag(string) (Droplets, error)
	Get(int) (*Droplet, error)
	Create(*godo.DropletCreateRequest, bool) (*Droplet, error)
	CreateMultiple(*godo.DropletMultiCreateRequest) (Droplets, error)
	Delete(int) error
	DeleteByTag(string) error
	Kernels(int) (Kernels, error)
	Snapshots(int) (Images, error)
	Backups(int) (Images, error)
	Actions(int) (Actions, error)
	Neighbors(int) (Droplets, error)
}

type dropletsService struct {
	client *godo.Client
}

var _ DropletsService = &dropletsService{}

// NewDropletsService builds a DropletsService instance.
func NewDropletsService(client *godo.Client) DropletsService {
	return &dropletsService{
		client: client,
	}
}

// func (ds *dropletsService) Iterator() <-chan *Droplet {
// 	ch := make(chan *Droplet)

//   go func() {

//   }
// }

func (ds *dropletsService) List() (Droplets, error) {
	f := func(opt *godo.ListOptions, out chan interface{}) (*godo.Response, error) {
		list, resp, err := ds.client.Droplets.List(opt)
		if err != nil {
			return nil, err
		}

		for _, d := range list {
			out <- d
		}

		return resp, nil
	}

	resp, err := PaginateResp(f)
	if err != nil {
		return nil, err
	}

	items := resp.([]interface{})
	list := make(Droplets, len(items))
	for i := range items {
		d := items[i].(godo.Droplet)
		list[i] = Droplet{Droplet: &d}
	}

	return list, nil
}

func (ds *dropletsService) ListByTag(tagName string) (Droplets, error) {
	f := func(opt *godo.ListOptions, out chan interface{}) (*godo.Response, error) {
		list, resp, err := ds.client.Droplets.ListByTag(tagName, opt)
		if err != nil {
			return nil, err
		}

		for _, d := range list {
			out <- d
		}

		return resp, nil
	}

	resp, err := PaginateResp(f)
	if err != nil {
		return nil, err
	}

	items := resp.([]interface{})
	list := make(Droplets, len(items))
	for i := range items {
		d := items[i].(godo.Droplet)
		list[i] = Droplet{Droplet: &d}
	}

	return list, nil
}

func (ds *dropletsService) Get(id int) (*Droplet, error) {
	d, _, err := ds.client.Droplets.Get(id)
	if err != nil {
		return nil, err
	}

	return &Droplet{Droplet: d}, nil
}

func (ds *dropletsService) Create(dcr *godo.DropletCreateRequest, wait bool) (*Droplet, error) {
	d, resp, err := ds.client.Droplets.Create(dcr)
	if err != nil {
		return nil, err
	}

	if wait {
		var action *godo.LinkAction
		for _, a := range resp.Links.Actions {
			if a.Rel == "create" {
				action = &a
				break
			}
		}

		if action != nil {
			_ = util.WaitForActive(ds.client, action.HREF)
			doDroplet, err := ds.Get(d.ID)
			if err != nil {
				return nil, err
			}
			d = doDroplet.Droplet
		}
	}

	return &Droplet{Droplet: d}, nil
}

func (ds *dropletsService) CreateMultiple(dmcr *godo.DropletMultiCreateRequest) (Droplets, error) {
	godoDroplets, _, err := ds.client.Droplets.CreateMultiple(dmcr)
	if err != nil {
		return nil, err
	}

	var droplets Droplets
	for _, d := range godoDroplets {
		droplets = append(droplets, Droplet{Droplet: &d})
	}

	return droplets, nil
}

func (ds *dropletsService) Delete(id int) error {
	_, err := ds.client.Droplets.Delete(id)
	return err
}

func (ds *dropletsService) DeleteByTag(tag string) error {
	_, err := ds.client.Droplets.DeleteByTag(tag)
	return err
}

func (ds *dropletsService) Kernels(id int) (Kernels, error) {
	f := func(opt *godo.ListOptions, out chan interface{}) (*godo.Response, error) {
		list, resp, err := ds.client.Droplets.Kernels(id, opt)
		if err != nil {
			return nil, err
		}

		for _, d := range list {
			out <- d
		}

		return resp, nil
	}

	resp, err := PaginateResp(f)
	if err != nil {
		return nil, err
	}

	items := resp.([]interface{})
	list := make(Kernels, len(items))
	for i := range items {
		d := items[i].(godo.Kernel)
		list[i] = Kernel{Kernel: &d}
	}

	return list, nil
}

func (ds *dropletsService) Snapshots(id int) (Images, error) {
	f := func(opt *godo.ListOptions, out chan interface{}) (*godo.Response, error) {
		list, resp, err := ds.client.Droplets.Snapshots(id, opt)
		if err != nil {
			return nil, err
		}

		for _, d := range list {
			out <- d
		}

		return resp, nil
	}

	resp, err := PaginateResp(f)
	if err != nil {
		return nil, err
	}

	items := resp.([]interface{})
	list := make(Images, len(items))
	for i := range items {
		d := items[i].(godo.Image)
		list[i] = Image{Image: &d}
	}

	return list, nil
}

func (ds *dropletsService) Backups(id int) (Images, error) {
	f := func(opt *godo.ListOptions, out chan interface{}) (*godo.Response, error) {
		list, resp, err := ds.client.Droplets.Backups(id, opt)
		if err != nil {
			return nil, err
		}

		for _, d := range list {
			out <- d
		}

		return resp, nil
	}

	resp, err := PaginateResp(f)
	if err != nil {
		return nil, err
	}

	items := resp.([]interface{})
	list := make(Images, len(items))
	for i := range items {
		d := items[i].(godo.Image)
		list[i] = Image{Image: &d}
	}

	return list, nil
}

func (ds *dropletsService) Actions(id int) (Actions, error) {
	f := func(opt *godo.ListOptions, out chan interface{}) (*godo.Response, error) {
		list, resp, err := ds.client.Droplets.Actions(id, opt)
		if err != nil {
			return nil, err
		}

		for _, d := range list {
			out <- d
		}

		return resp, nil
	}

	resp, err := PaginateResp(f)
	if err != nil {
		return nil, err
	}

	items := resp.([]interface{})
	list := make(Actions, len(items))
	for i := range items {
		d := items[i].(godo.Action)
		list[i] = Action{Action: &d}
	}
	return list, nil
}

func (ds *dropletsService) Neighbors(id int) (Droplets, error) {
	list, _, err := ds.client.Droplets.Neighbors(id)
	if err != nil {
		return nil, err
	}

	var droplets Droplets
	for _, d := range list {
		droplets = append(droplets, Droplet{Droplet: &d})
	}

	return droplets, nil
}
