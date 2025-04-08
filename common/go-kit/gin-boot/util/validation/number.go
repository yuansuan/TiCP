package validation

import (
	"fmt"
)

// Positive Positive
func Positive(num int) error {
	if num < 0 {
		return fmt.Errorf("ID(%v) < 0", num)
	}

	return nil
}

// ID ID
func ID(ids ...interface{}) error {
	for _, id := range ids {
		switch t := id.(type) {
		case int:
			return Positive(t)
		case []int:
			for _, i := range t {
				if err := Positive(i); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
