package main

import (
	"testing"
	"time"
)

func TestHalves(t *testing.T) {

	tables := []struct {
		initial int
		halves int
	} {
		{100, 7},
		{200, 8},
		{5, 2},
	}

	for _, table := range tables {
		caffHalves := Halves(table.initial)
		if caffHalves != table.halves {
			t.Errorf("Halvgorithm of %d was incorrect, got: %d, want: %d.",
				table.initial, caffHalves, table.halves)
		}
	}
}

func TestCalculateRemnants(t *testing.T) {
	dosage := 100
	bufSize := Halves(dosage)

	channel := make(chan Remnant, bufSize)
	go CalculateRemnants(Dose { Name: "Chai Tea Latte", Dosage: dosage,
		Time: time.Now().Format(time.RFC3339) }, channel)

	for remy := range channel {
		if remy.Amount >= dosage {
			t.Errorf("Remnant calculation of %d was incorrect, got %d",
				dosage, remy.Amount)
		}
	}
}
