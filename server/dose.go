package main

import "time"

// TODO: right now just using a strict average. Better to have a range? (upper, lower)
type Metabolizer struct {
        Name     string        `json:"name"`
        Onset    time.Duration `json:"onset"`
        Halflife time.Duration `json:"halflife"`
}
// niguding := Metabolizer { Name: "nicotine", Onset: time.Second*5, Halflife: time.Hour*2 }
// naipusheng := Metabolizer { Name: "naproxen", Onset: time.Minute*30, Halflife: time.Hour*13 }

type Remnant struct {
        Date   string `json:"time"`
        Amount int    `json:"remnant"`
}

type Dose struct {
        Id       int       `json:"id"` //TODO: reconcile this with ../dose.go which doesn't have it
        Name     string    `json:"name"`
        Dosage   int       `json:"dosage"`
        Time     string    `json:"time"`
        Remnants []Remnant `json:"remnants"`
}
