package util

import (
	"fmt"
	"sync"
)

// ChkErr ChkErr
func ChkErr(e error) {
	if e != nil {
		panic(e)
	}
}

// ErrorArrayBuilder build ErrorArray
type ErrorArrayBuilder ErrorArray

// ErrorArray is lots of error
type ErrorArray struct {
	l      sync.Mutex
	errors []error
}

// Error Error
func (p *ErrorArray) Error() string {
	return fmt.Sprintf("%v", p.errors)
}

// NewErrorArrayBuilder create ErrorArrayBuilder
func NewErrorArrayBuilder() *ErrorArrayBuilder {
	return &ErrorArrayBuilder{}
}

// Append : append a error to ErrorArray
func (p *ErrorArrayBuilder) Append(err error) *ErrorArrayBuilder {
	p.l.Lock()
	defer p.l.Unlock()

	if e, ok := err.(*ErrorArray); ok {
		p.errors = append(p.errors, e.errors...)
	} else {
		p.errors = append(p.errors, err)
	}
	return p
}

// Err : wrap as error
func (p *ErrorArrayBuilder) Err() error {
	p.l.Lock()
	defer p.l.Unlock()

	if len(p.errors) == 0 {
		return nil
	}
	return (*ErrorArray)(p)
}
