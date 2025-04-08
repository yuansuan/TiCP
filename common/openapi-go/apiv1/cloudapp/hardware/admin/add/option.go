package add

import (
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/hardware"
)

type Option func(req *hardware.AdminPostRequest) error

func (api API) Name(name string) Option {
	return func(req *hardware.AdminPostRequest) error {
		req.Name = &name
		return nil
	}
}

func (api API) Desc(desc string) Option {
	return func(req *hardware.AdminPostRequest) error {
		req.Desc = &desc
		return nil
	}
}

func (api API) InstanceType(instanceType string) Option {
	return func(req *hardware.AdminPostRequest) error {
		req.InstanceType = &instanceType
		return nil
	}
}

func (api API) InstanceFamily(instanceFamily string) Option {
	return func(req *hardware.AdminPostRequest) error {
		req.InstanceFamily = &instanceFamily
		return nil
	}
}

func (api API) Network(network int) Option {
	return func(req *hardware.AdminPostRequest) error {
		req.Network = &network
		return nil
	}
}

func (api API) Cpu(cpu int) Option {
	return func(req *hardware.AdminPostRequest) error {
		req.Cpu = &cpu
		return nil
	}
}

func (api API) Mem(mem int) Option {
	return func(req *hardware.AdminPostRequest) error {
		req.Mem = &mem
		return nil
	}
}

func (api API) Gpu(gpu int) Option {
	return func(req *hardware.AdminPostRequest) error {
		req.Gpu = &gpu
		return nil
	}
}

func (api API) GpuModel(gpuModel string) Option {
	return func(req *hardware.AdminPostRequest) error {
		req.GpuModel = &gpuModel
		return nil
	}
}

func (api API) CpuModel(cpuModel string) Option {
	return func(req *hardware.AdminPostRequest) error {
		req.CpuModel = &cpuModel
		return nil
	}
}

func (api API) Zone(zone string) Option {
	return func(req *hardware.AdminPostRequest) error {
		req.Zone = &zone
		return nil
	}
}
