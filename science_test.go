package main

import "testing"

func TestTrySpreadExpOrScienceInteger(t *testing.T) {
	type data struct {
		input string
		full  string
		isOk  bool
	}
	testData := []data{
		{"10^0", "1", true},
		{"10^1", "10", true},
		{"10^18", "1000000000000000000", true},
		{"2^10", "", false},
		{"10.1^2", "", false},
		{"10^-2", "", false},

		{"1e0", "1", true},
		{"1E1", "10", true},
		{"1e18", "1000000000000000000", true},
		{"1.2345e4", "12345", true},
		{"2.3e4", "23000", true},
		{"1.23456e4", "", false},
		{".23e4", "", false},
		{"12e4", "", false},
	}

	for _, d := range testData {
		f, ok := TrySpreadExpOrScienceInteger(d.input)
		if ok != d.isOk || f != d.full {
			t.Errorf("%s want: %q %v ; got: %q %v\n", d.input, d.full, d.isOk, f, ok)
		}
	}
}
