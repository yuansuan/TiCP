package publish

import "fmt"

type Status int

const (
	Default     = Status(0)
	Published   = Status(1)
	Unpublished = Status(2)
)

type UnsupportedStatusError struct {
	Value int
}

func (s UnsupportedStatusError) Error() string {
	return fmt.Sprintf("Unsupported status: %d", s.Value)
}

func ToStatus(value int) (Status, error) {
	switch Status(value) {
	case Default:
		return Default, nil
	case Published:
		return Published, nil
	case Unpublished:
		return Unpublished, nil
	default:
		return -1, UnsupportedStatusError{Value: value}
	}
}

func (s *Status) Published() bool {
	return *s == Published
}

func (s *Status) GetValue() int {
	return int(*s)
}
