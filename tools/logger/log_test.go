package logger

import "testing"

func TestErr(t *testing.T) {
	Err("123", "test error")
}
