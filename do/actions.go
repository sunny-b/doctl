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

import "github.com/digitalocean/godo"

//Action is a wrapper for godo.Action
type Action struct {
	*godo.Action
}

// Actions is a slice of Action.
type Actions []Action

// ActionsService is an interface for interacting with DigitalOcean's action api.
type ActionsService interface {
	List() (Actions, error)
	Get(int) (*Action, error)
}

type actionsService struct {
	client *godo.Client
}

var _ ActionsService = &actionsService{}

// NewActionsService builds an ActionsService instance.
func NewActionsService(godoClient *godo.Client) ActionsService {
	return &actionsService{
		client: godoClient,
	}
}

func (as *actionsService) List() (Actions, error) {
	f := func(opt *godo.ListOptions, out chan interface{}) (*godo.Response, error) {
		list, resp, err := as.client.Actions.List(opt)
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

func (as *actionsService) Get(id int) (*Action, error) {
	a, _, err := as.client.Actions.Get(id)
	if err != nil {
		return nil, err
	}

	return &Action{Action: a}, nil
}
