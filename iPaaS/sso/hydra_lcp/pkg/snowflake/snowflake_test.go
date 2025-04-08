package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSnowflake(t *testing.T) {
	// 0 convert to empty string
	zeroID := ParseInt64(int64(0))
	assert := require.New(t)
	assert.Equal("", zeroID.String())
}
