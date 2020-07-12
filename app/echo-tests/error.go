package main

import (
	"fmt"
)

func (r ServiceError) Error() string {
	return fmt.Sprintf("status %d: err %v :msg %v", r.Code, r.Err, r.Message)
}

// ServiceError handles custom errors
type ServiceError struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Err     error  `json:"-"`
}
