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

// Image is a werapper for godo.Image
type Image struct {
	*godo.Image
}

// Images is a slice of Droplet.
type Images []Image

// ImagesService is the godo ImagesService interface.
type ImagesService interface {
	List(public bool) (Images, error)
	ListDistribution(public bool) (Images, error)
	ListApplication(public bool) (Images, error)
	ListUser(public bool) (Images, error)
	GetByID(id int) (*Image, error)
	GetBySlug(slug string) (*Image, error)
	Update(id int, iur *godo.ImageUpdateRequest) (*Image, error)
	Delete(id int) error
}

type imagesService struct {
	client *godo.Client
}

var _ ImagesService = &imagesService{}

// NewImagesService builds an instance of ImagesService.
func NewImagesService(client *godo.Client) ImagesService {
	return &imagesService{
		client: client,
	}
}

func (is *imagesService) List(public bool) (Images, error) {
	return is.listImages(is.client.Images.List, public)
}

func (is *imagesService) ListDistribution(public bool) (Images, error) {
	return is.listImages(is.client.Images.ListDistribution, public)
}

func (is *imagesService) ListApplication(public bool) (Images, error) {
	return is.listImages(is.client.Images.ListApplication, public)
}

func (is *imagesService) ListUser(public bool) (Images, error) {
	return is.listImages(is.client.Images.ListUser, public)
}

func (is *imagesService) GetByID(id int) (*Image, error) {
	i, _, err := is.client.Images.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &Image{Image: i}, nil
}

func (is *imagesService) GetBySlug(slug string) (*Image, error) {
	i, _, err := is.client.Images.GetBySlug(slug)
	if err != nil {
		return nil, err
	}

	return &Image{Image: i}, nil
}

func (is *imagesService) Update(id int, iur *godo.ImageUpdateRequest) (*Image, error) {
	i, _, err := is.client.Images.Update(id, iur)
	if err != nil {
		return nil, err
	}

	return &Image{Image: i}, nil
}

func (is *imagesService) Delete(id int) error {
	_, err := is.client.Images.Delete(id)
	return err
}

type listFn func(*godo.ListOptions) ([]godo.Image, *godo.Response, error)

func (is *imagesService) listImages(lFn listFn, public bool) (Images, error) {
	fn := func(opt *godo.ListOptions, out chan interface{}) (*godo.Response, error) {
		list, resp, err := lFn(opt)
		if err != nil {
			return nil, err
		}

		for _, d := range list {
			out <- d
		}

		return resp, nil
	}

	resp, err := PaginateResp(fn)
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
