package validation

import (
	"fmt"
	"regexp"
)

// Name validates name.
func Name(name string) error {
	match, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_]*$`, name)
	if !match {
		return fmt.Errorf("Name '%v' should be made up of a-zA-z0-9_, start with a-zA-Z", name)
	}

	return nil
}

// Email validates email.
func Email(email string) error {
	match, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`, email)
	if !match {
		return fmt.Errorf("Email \"%v\" format wrong", email)
	}

	return nil
}
