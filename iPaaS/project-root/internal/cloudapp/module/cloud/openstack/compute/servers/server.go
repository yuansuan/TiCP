package servers

import (
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/extendedstatus"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

type Server struct {
	servers.Server
	extendedstatus.ServerExtendedStatusExt
}
