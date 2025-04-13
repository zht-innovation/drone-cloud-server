package tests

import (
	"runtime"
	"strings"
	"testing"
)

func getCallerInfo() (string, int) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return "unknown", 0
	}

	startIdx := strings.Index(file, "/cloud")
	return file[startIdx:], line
}

func TestLog(t *testing.T) {
	file, line := getCallerInfo()
	if file != "/cloud/tests/logger_test.go" {
		t.Errorf("file: '%s' is not equal to '/cloud/tests/logger_test.go'", file)
	}

	if line != 20 {
		t.Errorf("linenum: '%d' is not equal to 20", line)
	}
}
