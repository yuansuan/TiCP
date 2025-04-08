package guacamole

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildProtocol(t *testing.T) {
	assert.Equal(t, invalidProtocolMsg, BuildProtocol("invalid"))
	assert.Equal(t, "10.storage_id,4.abcd;", BuildProtocol("storage_id", "abcd"))
}
