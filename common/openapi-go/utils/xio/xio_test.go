package xio

import (
	"crypto/rand"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDupReader(t *testing.T) {
	r := io.NopCloser(io.LimitReader(rand.Reader, 1024))
	if r1, r2, err := DupReader(r); assert.NoError(t, err) {
		if data1, err := ioutil.ReadAll(r1); assert.NoError(t, err) {
			if data2, err := ioutil.ReadAll(r2); assert.NoError(t, err) {
				assert.Equal(t, data1, data2)
			}
		}

		assert.NoError(t, r1.Close())
		assert.NoError(t, r2.Close())
	}
}
