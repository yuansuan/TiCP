package v20230530

import (
	"time"
)

type Hardware struct {
	HardwareId     string `json:"HardwareId,omitempty"`
	Zone           string `json:"Zone,omitempty" `
	Name           string `json:"Name,omitempty" `
	Desc           string `json:"Desc,omitempty" `
	InstanceType   string `json:"InstanceType,omitempty" `
	InstanceFamily string `json:"InstanceFamily,omitempty" `
	Network        int    `json:"Network,omitempty" `
	Cpu            int    `json:"Cpu,omitempty" `
	CpuModel       string `json:"CpuModel,omitempty" `
	Mem            int    `json:"Mem,omitempty" `
	Gpu            int    `json:"Gpu,omitempty" `
	GpuModel       string `json:"GpuModel,omitempty" `
}

type Software struct {
	SoftwareId string       `json:"SoftwareId,omitempty"`
	Zone       string       `json:"Zone,omitempty" `
	Name       string       `json:"Name,omitempty" `
	Desc       string       `json:"Desc,omitempty" `
	Icon       string       `json:"Icon,omitempty" `
	Platform   string       `json:"Platform,omitempty" `
	ImageId    string       `json:"ImageId,omitempty" `
	InitScript string       `json:"InitScript,omitempty" `
	GpuDesired *bool        `json:"GpuDesired,omitempty" `
	RemoteApps []*RemoteApp `json:"RemoteApps,omitempty"`
}

type RemoteApp struct {
	Id         string  `json:"Id,omitempty"`
	SoftwareId *string `json:"SoftwareId,omitempty"`
	Desc       *string `json:"Desc,omitempty"`
	Name       *string `json:"Name,omitempty"`
	Dir        *string `json:"Dir,omitempty"`
	Args       *string `json:"Args,omitempty"`
	Logo       *string `json:"Logo,omitempty"`
	DisableGfx *bool   `json:"DisableGfx,omitempty"`
	LoginUser  *string `json:"LoginUser,omitempty"`
}

type Session struct {
	Id              string       `json:"Id,omitempty"`
	UserId          string       `json:"UserId,omitempty"`
	Zone            string       `json:"Zone,omitempty"`
	Status          string       `json:"Status,omitempty"`
	StreamUrl       string       `json:"StreamUrl,omitempty"`
	CreateTime      *time.Time   `json:"CreateTime,omitempty"`
	StartTime       *time.Time   `json:"StartTime,omitempty"`
	EndTime         *time.Time   `json:"EndTime,omitempty"`
	MachinePassword string       `json:"MachinePassword,omitempty"`
	ExitReason      string       `json:"ExitReason,omitempty"`
	Software        *Software    `json:"Software,omitempty"`
	Hardware        *Hardware    `json:"Hardware,omitempty"`
	RemoteApps      []*RemoteApp `json:"RemoteApps,omitempty"`
}

type SessionStatus string

const (
	SessionPending     SessionStatus = "PENDING"
	SessionStarting    SessionStatus = "STARTING"
	SessionStarted     SessionStatus = "STARTED"
	SessionClosing     SessionStatus = "CLOSING"
	SessionClosed      SessionStatus = "CLOSED"
	SessionPoweringOff SessionStatus = "POWERING OFF"
	SessionPowerOff    SessionStatus = "POWER OFF"
	SessionPoweringOn  SessionStatus = "POWERING ON"
	SessionRebooting   SessionStatus = "REBOOTING"
)

var SessionStatusMap = map[SessionStatus]string{
	SessionPending:     "PENDING",
	SessionStarting:    "STARTING",
	SessionStarted:     "STARTED",
	SessionClosing:     "CLOSING",
	SessionClosed:      "CLOSED",
	SessionPoweringOff: "POWERING OFF",
	SessionPowerOff:    "POWER OFF",
	SessionPoweringOn:  "POWERING ON",
	SessionRebooting:   "REBOOTING",
}

func (s SessionStatus) String() string {
	return SessionStatusMap[s]
}
