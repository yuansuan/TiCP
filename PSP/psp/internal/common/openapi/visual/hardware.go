package openapivisual

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/hardware"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	visualdto "github.com/yuansuan/ticp/PSP/psp/internal/visual/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

func AddHardware(api *openapi.OpenAPI, hardware *visualdto.Hardware, zone string) (*hardware.AdminPostResponse, error) {
	if strutil.IsEmpty(hardware.Name) {
		return nil, ErrHardwareNameEmpty
	}
	if strutil.IsEmpty(hardware.InstanceType) {
		return nil, ErrHardwareInstanceTypeEmpty
	}
	if strutil.IsEmpty(hardware.InstanceFamily) {
		return nil, ErrHardwareInstanceFamilyEmpty
	}
	if strutil.IsEmpty(zone) {
		return nil, ErrHardwareZoneEmpty
	}

	client := api.Client.CloudApp.Hardware.Admin
	response, err := api.Client.CloudApp.Hardware.Admin.Add(
		client.Add.Name(hardware.Name),
		client.Add.Desc(hardware.Desc),
		client.Add.Network(hardware.Network),
		client.Add.Cpu(hardware.CPU),
		client.Add.Mem(hardware.Mem),
		client.Add.Gpu(hardware.GPU),
		client.Add.CpuModel(hardware.CPUModel),
		client.Add.GpuModel(hardware.GPUModel),
		client.Add.InstanceType(hardware.InstanceType),
		client.Add.InstanceFamily(hardware.InstanceFamily),
		client.Add.Zone(zone),
	)

	if response != nil {
		logging.Default().Debugf("openapi add hardware request id: [%v], req: [%+v], zone: [%v]", response.RequestID, hardware, zone)
	}

	return response, err
}

func UpdateHardware(api *openapi.OpenAPI, hardwareID string, hardware *visualdto.Hardware, zone string) (*hardware.AdminPutResponse, error) {
	if strutil.IsEmpty(hardwareID) {
		return nil, ErrHardwareIDEmpty
	}
	if strutil.IsEmpty(hardware.Name) {
		return nil, ErrHardwareNameEmpty
	}
	if strutil.IsEmpty(hardware.InstanceType) {
		return nil, ErrHardwareInstanceTypeEmpty
	}
	if strutil.IsEmpty(hardware.InstanceFamily) {
		return nil, ErrHardwareInstanceFamilyEmpty
	}
	if strutil.IsEmpty(zone) {
		return nil, ErrHardwareZoneEmpty
	}

	client := api.Client.CloudApp.Hardware.Admin
	response, err := api.Client.CloudApp.Hardware.Admin.Put(
		client.Put.Id(hardwareID),
		client.Put.Name(hardware.Name),
		client.Put.Desc(hardware.Desc),
		client.Put.Network(hardware.Network),
		client.Put.Cpu(hardware.CPU),
		client.Put.Mem(hardware.Mem),
		client.Put.Gpu(hardware.GPU),
		client.Put.CpuModel(hardware.CPUModel),
		client.Put.GpuModel(hardware.GPUModel),
		client.Put.InstanceType(hardware.InstanceType),
		client.Put.InstanceFamily(hardware.InstanceFamily),
		client.Put.Zone(zone),
	)

	if response != nil {
		logging.Default().Debugf("openapi update hardware request id: [%v], hardwareID: [%v], req: [%+v], "+
			"zone: [%v]", response.RequestID, hardwareID, hardware, zone)
	}

	return response, err
}

func DeleteHardware(api *openapi.OpenAPI, hardwareID string) (*hardware.AdminDeleteResponse, error) {
	if strutil.IsEmpty(hardwareID) {
		return nil, ErrHardwareIDEmpty
	}

	client := api.Client.CloudApp.Hardware.Admin
	response, err := api.Client.CloudApp.Hardware.Admin.Delete(
		client.Delete.Id(hardwareID),
	)

	if response != nil {
		logging.Default().Debugf("openapi delete hardware request id: [%v], hardwareID: [%v]", response.RequestID, hardwareID)
	}

	return response, err
}
