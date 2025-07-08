package tests

import (
	"testing"
)

func validateMacFormat(mac string) bool {
	if len(mac) != 17 {
		return false
	}
	for i, c := range mac {
		if (i+1)%3 == 0 {
			if c != ':' && c != '-' {
				return false
			}
		} else {
			if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F') {
				return false
			}
		}
	}

	return true
}

func TestValidateMacFormat(t *testing.T) {
	testcases := []struct {
		mac      string
		expected bool
	}{
		{"00:11:22:33:44:55", true},
		{"00-11-22-33-44-55", true},
		{"001122334455", false},
		{"00:11:22:33:44", false},
		{"00:11:22:33:44:55:66", false},
		{"00:11:22:33:44:G5", false},
		{"0A-7B-8F-9F-98-1B", true},
	}

	for _, test := range testcases {
		result := validateMacFormat(test.mac)
		if result != test.expected {
			t.Errorf("validateMacFormat(%s) = %v; expected %v", test.mac, result, test.expected)
		}
	}
}
