package main

import (
	"fmt"
	"runtime"
	"testing"
)

func TestGetFilename(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	fmt.Println("Current test file: " + filename)
}
