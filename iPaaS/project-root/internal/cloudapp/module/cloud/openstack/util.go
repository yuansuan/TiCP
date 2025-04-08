package openstack

import (
	"strings"
)

func stringSliceToString(strSlice []string) string {
	return strings.Join(strSlice, ",")
}
