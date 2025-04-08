package cloud

import (
	"fmt"
	"strings"
)

type ScriptParams struct {
	SignalHost          string          `json:"signal_host,omitempty"`
	RoomId              string          `json:"room_id,omitempty"`
	ShareUsername       StringList      `json:"share_username,omitempty"`
	SharePassword       StringList      `json:"share_password,omitempty"`
	ShareMountPaths     ShareMountPaths `json:"share_mount_paths,omitempty"` // encoded
	LoginPassword       string          `json:"login_password,omitempty"`
	RemoteAppUserPasses string          `json:"remote_app_user_passes,omitempty"`
}

func (sp *ScriptParams) UpdateByMount(shareUsername, sharePassword, mountSrc, mountPoint string) *ScriptParams {
	sp.ShareUsername = sp.ShareUsername.Append(shareUsername)
	sp.SharePassword = sp.SharePassword.Append(sharePassword)
	sp.ShareMountPaths = sp.ShareMountPaths.Append(mountSrc, mountPoint)

	return sp
}

func (sp *ScriptParams) UpdateByUMount(mountPoint string) *ScriptParams {
	// find index need delete
	deleteIndex := sp.ShareMountPaths.GetMountPointIndex(mountPoint)
	if deleteIndex == -1 {
		return sp
	}

	sp.ShareUsername = sp.ShareUsername.DelByIndex(deleteIndex)
	sp.SharePassword = sp.SharePassword.DelByIndex(deleteIndex)
	sp.ShareMountPaths = sp.ShareMountPaths.DelByIndex(deleteIndex)

	return sp
}

type StringList string

func (l StringList) String() string {
	return string(l)
}

func (l StringList) Append(s string) StringList {
	if l == "" {
		return StringList(s)
	}

	return StringList(fmt.Sprintf("%s,%s", l, s))
}

func (l StringList) toStringSlice() []string {
	return strings.Split(l.String(), ",")
}

func (l StringList) DelByIndex(index int) StringList {
	arr := make([]string, 0)
	for i, v := range l.toStringSlice() {
		if i != index {
			arr = append(arr, v)
		}
	}

	return StringList(strings.Join(arr, ","))
}

type ShareMountPaths string

func (smp ShareMountPaths) String() string {
	return string(smp)
}

func (smp ShareMountPaths) Append(mountSrc, mountPoint string) ShareMountPaths {
	return ShareMountPaths(StringList(smp).Append(fmt.Sprintf("%s=%s", mountSrc, mountPoint)))
}

func (smp ShareMountPaths) toStringSlice() []string {
	return StringList(smp).toStringSlice()
}

func (smp ShareMountPaths) IsMountPointExist(mountPoint string) bool {
	for _, s := range smp.toStringSlice() {
		kv := strings.Split(s, "=")
		if len(kv) != 2 {
			continue
		}

		if mountPoint == kv[1] {
			return true
		}
	}

	return false
}

func (smp ShareMountPaths) DelByIndex(index int) ShareMountPaths {
	arr := make([]string, 0)
	for i, v := range smp.toStringSlice() {
		if i != index {
			arr = append(arr, v)
		}
	}

	return ShareMountPaths(strings.Join(arr, ","))
}

func (smp ShareMountPaths) GetMountPointIndex(mountPoint string) int {
	for i, s := range smp.toStringSlice() {
		kv := strings.Split(s, "=")
		if len(kv) != 2 {
			continue
		}

		if mountPoint == kv[1] {
			return i
		}
	}

	return -1
}
