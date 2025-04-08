package jsonmap

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp_dirnfs/module/xio"
)

func TestLoad(t *testing.T) {
	fp, err := os.CreateTemp(os.TempDir(), "*")
	require.NoError(t, err)
	defer func() { _ = os.Remove(fp.Name()) }()
	defer func() { _ = fp.Close() }()

	m, err := Load(fp)
	require.NoError(t, err)

	if m.Set("aa", 123); assert.Equal(t, 123, m.Get("aa")) {
		var buf bytes.Buffer
		if err = xio.WriteFrom(&buf, fp, 0); assert.NoError(t, err) {
			t.Logf("data: %s", buf.Bytes())
		}
	}
}
