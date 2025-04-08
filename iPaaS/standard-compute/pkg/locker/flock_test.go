package locker

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/fsutil/filemode"
)

func TestFileLocker_Lock_WithBlocking(t *testing.T) {
	random := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, random); assert.NoError(t, err) {
		temp := filepath.Join(os.TempDir(), hex.EncodeToString(random))
		open := func() (*os.File, func(), error) {
			f, err := os.OpenFile(temp, os.O_CREATE|os.O_RDWR, filemode.RegularFile)
			return f, func() {
				_ = os.RemoveAll(temp)
			}, err
		}

		var locked bool

		l1 := NewFileLocker(open)
		if ok, err := l1.Lock(FastFail); assert.True(t, ok) && assert.NoError(t, err) {
			go func() {
				l2 := NewFileLocker(open)
				if ok, err := l2.Lock(UntilAcquire); assert.True(t, ok) && assert.NoError(t, err) {
					defer func() { _ = l2.Unlock() }()
					assert.True(t, locked)
				}
			}()

			time.Sleep(time.Second)
			locked = true
			_ = l1.Unlock()
		}
	}
}

func TestFileLocker_Lock_WithTimeout(t *testing.T) {
	random := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, random); assert.NoError(t, err) {
		temp := filepath.Join(os.TempDir(), hex.EncodeToString(random))
		open := func() (*os.File, func(), error) {
			f, err := os.OpenFile(temp, os.O_CREATE|os.O_RDWR, filemode.RegularFile)
			return f, func() {
				_ = os.RemoveAll(temp)
			}, err
		}

		l1 := NewFileLocker(open)
		if ok, err := l1.Lock(FastFail); assert.True(t, ok) && assert.NoError(t, err) {
			go func() {
				l2 := NewFileLocker(open)
				ok, err := l2.Lock(time.Second / 2)
				assert.False(t, ok)
				assert.ErrorIs(t, err, ErrLockTimeout)
			}()

			time.Sleep(time.Second)
			_ = l1.Unlock()
		}
	}
}
