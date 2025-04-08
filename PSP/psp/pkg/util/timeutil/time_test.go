package timeutil

import (
	"testing"
	"time"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
)

func TestDefaultDateTime(t *testing.T) {
	if DefaultDateTime.Format(time.RFC3339) != common.DefaultEmptyTime {
		t.Fatal("default time validate failed")
	}
}
