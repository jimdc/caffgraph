package main

import "testing"

func TestCalculateRemnants(t *testing.T) {

	bufSize := Halves(100)
	if bufSize != 7 {
		t.Errorf("bufSize was incorrect, got: %d, want: %d.", bufSize, 7)
	}

	//channel := make(chan Remnant, bufSize)
	//go CalculateRemnants(time.Now(), 100, channel)
	//for remy := range channel { }

	//Next: do test tables https://blog.alexellis.io/golang-writing-unit-tests/
}
