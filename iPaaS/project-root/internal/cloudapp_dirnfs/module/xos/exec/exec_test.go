package exec

import "testing"

func TestNew(t *testing.T) {
	New("echo", String("abc"), String("def"), PlaceHolder(1))
}
