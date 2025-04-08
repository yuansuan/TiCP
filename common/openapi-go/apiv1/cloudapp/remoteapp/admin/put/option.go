package patch

import "github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/remoteapp"

type Option func(req *remoteapp.AdminPutRequest) error

func (api API) Id(remoteAppId string) Option {
	return func(req *remoteapp.AdminPutRequest) error {
		req.RemoteAppId = &remoteAppId
		return nil
	}
}

func (api API) SoftwareId(softwareId string) Option {
	return func(req *remoteapp.AdminPutRequest) error {
		req.SoftwareId = &softwareId
		return nil
	}
}

func (api API) Name(name string) Option {
	return func(req *remoteapp.AdminPutRequest) error {
		req.Name = &name
		return nil
	}
}

func (api API) Desc(desc string) Option {
	return func(req *remoteapp.AdminPutRequest) error {
		req.Desc = &desc
		return nil
	}
}

func (api API) Dir(dir string) Option {
	return func(req *remoteapp.AdminPutRequest) error {
		req.Dir = &dir
		return nil
	}
}

func (api API) Args(args string) Option {
	return func(req *remoteapp.AdminPutRequest) error {
		req.Args = &args
		return nil
	}
}

func (api API) Logo(logo string) Option {
	return func(req *remoteapp.AdminPutRequest) error {
		req.Logo = &logo
		return nil
	}
}

func (api API) DisableGfx(disableGfx bool) Option {
	return func(req *remoteapp.AdminPutRequest) error {
		req.DisableGfx = &disableGfx
		return nil
	}
}
