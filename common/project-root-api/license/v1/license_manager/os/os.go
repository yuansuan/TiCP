package os

import (
	"fmt"
)

type OS int

const (
	Default = OS(0)
	Linux   = OS(1)
	Windows = OS(2)
)

type UnsupportedOSError struct {
	Value int
}

func (e UnsupportedOSError) Error() string {
	return fmt.Sprintf("Unsupported OS: %d", e.Value)
}

func ToOS(value int) (OS, error) {
	switch OS(value) {
	case Default:
		return Default, nil
	case Linux:
		return Linux, nil
	case Windows:
		return Windows, nil
	default:
		return -1, UnsupportedOSError{Value: value}
	}
}

func (t *OS) GetValue() int {
	return int(*t)
}
