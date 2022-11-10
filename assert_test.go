package grss

import (
	"fmt"
	"runtime/debug"
	"testing"
)

// I believe testify is better, but I want to keep zero-dependency
func assert(t *testing.T, b bool, args ...interface{}) {
	if !b {
		for i := range args {
			switch args[i].(type) {
			case error, nil:
				continue
			default:
				args[i] = fmt.Sprintf("%+v\n", args[i])
			}
		}
		fmt.Println(string(debug.Stack()))
		t.Error(args...)
	}
}
