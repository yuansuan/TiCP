package util

import (
	"fmt"
	"strings"

	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/collection"
)

// Closer Closer
type Closer interface {
	Close() error
}

// Close is a convenience function to close a object that has a Close() method, ignoring any errors
// Used to satisfy errcheck lint
func Close(c Closer) {
	if err := c.Close(); err != nil {
		logging.Default().Warnf("failed to close %v: %v", c, err)
	}
}

// CheckEnvKey checks the user key of environment variable whether or not is conflict with system
// Ignore the user one if conflict
func CheckEnvKey(key string) bool {
	var blackList = []string{"MANPATH", "HOSTNAME", "SHELL", "TERM", "USER", "PATH", "HOME", "LOGNAME", "_", "PWD"}

	if collection.Contain(blackList, key) {
		return false
	}

	return true
}

// ConvertBoolToInt ConvertBoolToInt
func ConvertBoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// ConvertIntToBool ConvertIntToBool
func ConvertIntToBool(i int) bool {
	return i == 1
}

// SCPrefix SCPrefix
const SCPrefix = "sc_"

// ConvertPeerID2SCId Convert peerID to SC id
// The convert rule is to remove the prefix "sc_" (if exists) of the peerID to get the scID
func ConvertPeerID2SCId(peerID string) string {
	scID := peerID // set default value
	if strings.HasPrefix(peerID, SCPrefix) {
		scID = strings.Replace(peerID, SCPrefix, "", 1)
	}
	return scID
}

// ConvertSCId2PeerID Convert SC id to peerID
// The convert rule is to add the prefix "sc_" (if exists) of the peerID to the scID
func ConvertSCId2PeerID(scID string) string {
	peerID := scID // set default value
	if !strings.HasPrefix(scID, SCPrefix) {
		peerID = fmt.Sprintf("%v%v", SCPrefix, scID)
	}
	return peerID
}
