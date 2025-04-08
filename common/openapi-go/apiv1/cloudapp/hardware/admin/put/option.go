package put

import "github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/hardware"

type Option func(req *hardware.AdminPutRequest) error

func (api API) Id(id string) Option {
	return func(req *hardware.AdminPutRequest) error {
		req.HardwareId = &id
		return nil
	}
}

func (api API) Name(name string) Option {
	return func(req *hardware.AdminPutRequest) error {
		req.Name = &name
		return nil
	}
}

func (api API) Desc(desc string) Option {
	return func(req *hardware.AdminPutRequest) error {
		req.Desc = &desc
		return nil
	}
}

func (api API) InstanceType(instanceType string) Option {
	return func(req *hardware.AdminPutRequest) error {
		req.InstanceType = &instanceType
		return nil
	}
}

func (api API) InstanceFamily(instanceFamily string) Option {
	return func(req *hardware.AdminPutRequest) error {
		req.InstanceFamily = &instanceFamily
		return nil
	}
}

func (api API) Network(network int) Option {
	return func(req *hardware.AdminPutRequest) error {
		req.Network = &network
		return nil
	}
}

func (api API) Cpu(cpu int) Option {
	return func(req *hardware.AdminPutRequest) error {
		req.Cpu = &cpu
		return nil
	}
}

func (api API) Mem(mem int) Option {
	return func(req *hardware.AdminPutRequest) error {
		req.Mem = &mem
		return nil
	}
}

func (api API) Gpu(gpu int) Option {
	return func(req *hardware.AdminPutRequest) error {
		req.Gpu = &gpu
		return nil
	}
}

func (api API) GpuModel(gpuModel string) Option {
	return func(req *hardware.AdminPutRequest) error {
		req.GpuModel = &gpuModel
		return nil
	}
}

func (api API) CpuModel(cpuModel string) Option {
	return func(req *hardware.AdminPutRequest) error {
		req.CpuModel = &cpuModel
		return nil
	}
}

func (api API) Zone(zone string) Option {
	return func(req *hardware.AdminPutRequest) error {
		req.Zone = &zone
		return nil
	}
}
