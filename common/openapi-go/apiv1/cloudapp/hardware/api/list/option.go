package list

import "github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/hardware"

type Option func(req *hardware.APIListRequest) error

func (api API) Name(name string) Option {
	return func(req *hardware.APIListRequest) error {
		req.Name = &name
		return nil
	}
}

func (api API) Cpu(cpu int) Option {
	return func(req *hardware.APIListRequest) error {
		req.Cpu = &cpu
		return nil
	}
}

func (api API) Mem(mem int) Option {
	return func(req *hardware.APIListRequest) error {
		req.Mem = &mem
		return nil
	}
}

func (api API) Gpu(gpu int) Option {
	return func(req *hardware.APIListRequest) error {
		req.Gpu = &gpu
		return nil
	}
}

func (api API) Zone(zone string) Option {
	return func(req *hardware.APIListRequest) error {
		req.Zone = &zone
		return nil
	}
}

func (api API) PageOffset(pageOffset int) Option {
	return func(req *hardware.APIListRequest) error {
		req.PageOffset = &pageOffset
		return nil
	}
}

func (api API) PageSize(pageSize int) Option {
	return func(req *hardware.APIListRequest) error {
		req.PageSize = &pageSize
		return nil
	}
}
