package do

import "github.com/digitalocean/godo"

// Volume is a wrapper for godo.Volume.
type Volume struct {
	*godo.Volume
}

// Volumes is a slice of Volume.
type Volumes []Volume

// VolumesService is an interface for interacting with DigitalOcean's account api.
type VolumesService interface {
	List() ([]Volume, error)
	CreateVolume(*godo.VolumeCreateRequest) (*Volume, error)
	DeleteVolume(string) error
	Get(string) (*Volume, error)
}

type volumesService struct {
	client *godo.Client
}

var _ VolumesService = &volumesService{}

// NewVolumesService builds an NewVolumesService instance.
func NewVolumesService(godoClient *godo.Client) VolumesService {
	return &volumesService{
		client: godoClient,
	}

}

func (a *volumesService) List() ([]Volume, error) {
	f := func(opt *godo.ListOptions, out chan interface{}) (*godo.Response, error) {
		list, resp, err := a.client.Storage.ListVolumes(opt)
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
	list := make(Volumes, len(items))
	for i := range items {
		d := items[i].(godo.Volume)
		list[i] = Volume{Volume: &d}
	}

	return list, nil
}

func (a *volumesService) CreateVolume(r *godo.VolumeCreateRequest) (*Volume, error) {
	al, _, err := a.client.Storage.CreateVolume(r)
	if err != nil {
		return nil, err

	}
	return &Volume{Volume: al}, nil

}

func (a *volumesService) DeleteVolume(id string) error {
	_, err := a.client.Storage.DeleteVolume(id)
	return err
}

func (a *volumesService) Get(id string) (*Volume, error) {
	d, _, err := a.client.Storage.GetVolume(id)
	if err != nil {
		return nil, err

	}

	return &Volume{Volume: d}, nil

}
