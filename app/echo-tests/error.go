package main

import (
	"fmt"
)

func (r ServiceError) Error() string {
	return fmt.Sprintf("status %d: err %v :msg %v", r.code, r.Err, r.Message)
}

// ServiceError handles custom errors
type ServiceError struct {
	code int
	Message    string
	Err        error
}
