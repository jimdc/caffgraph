package main

import (
	"io/ioutil"
	"flag"
	"fmt"
	"math"
	"os"
	"time"
	"encoding/json"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// math.Round doesn't round up, so compensate when allocating
func Halves(mg int) (timesToDivide int) {
	fmg := float64(mg)
	nDivides := math.Round(math.Log2(fmg))
	return int(nDivides)
}

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
        Name     string    `json:"name"`
        Dosage   int       `json:"dosage"`
        Time     string    `json:"time"`
        Remnants []Remnant `json:"remnants"`
}

/*
type Temple struct {
	Offerings []Dose           `json:"offerings"`
	Metabolizers []Metabolizer `json:"metabolizers"`
}
*/

func Metabolize(damage Dose, instruction Metabolizer, Results chan Remnant) {

	t, err := time.Parse(time.RFC3339, damage.Time)
	check(err)
	ts := damage.Time

	for remaining := damage.Dosage / 2; remaining >= 1; remaining /= 2 {
		Results <- Remnant{Date: ts, Amount: remaining}

		t = t.Add(instruction.Halflife)
		ts = t.Format(time.RFC3339)
	}

	close(Results)
}

const Filename string = "caff.json"

// This function name will make more sense when it isn't only []Dose but also include Alim info.
func LoadTemple() (doses []Dose) {
	jsonFile, err := os.Open(Filename)
	check(err)

	byteValue, err := ioutil.ReadAll(jsonFile)
	check(err)

	jsonFile.Close()

	var ds []Dose
	err = json.Unmarshal(byteValue, &ds)
	check(err)

	return ds
}

func main() {
	flagWrite := flag.Bool("write", false, "write output to " + Filename)
	flagTime := flag.String("time", "now", "time of dosage, e.g. 2018-04-16T17:22:40Z")
	flagDosage := flag.Int("dosage", 100, "caffeine dosage in mg")
	flagName := flag.String("name", "Americano", "name of caffeine product")
	flagRead := flag.Bool("read", false, "read json in " + Filename)
	flag.Parse()
        // TODO: see if it's possible to hide some value in struct from json parser? All this conversion seems ineff
        t := time.Now()
        var err error
        if *flagTime != "now" {
                t, err = time.Parse(time.RFC3339, *flagTime)
                check(err)
        }
        ts := t.Format(time.RFC3339)
        inputDose := Dose{Name: *flagName, Dosage: *flagDosage, Time: ts, Remnants: nil}

	//TODO: integrate with write, to not have duplicate arrays and stuff
	if *flagRead == true {
		doses := LoadTemple()
		fmt.Printf("%v\n", doses)
		os.Exit(0)
	}

	var f *os.File
	if *flagWrite == true {
		f, err = os.OpenFile(Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		check(err)
		defer f.Close()
	}

	bufSize := Halves(*flagDosage)
	channel := make(chan Remnant, bufSize)

	caffeine := Metabolizer { Name: "caffeine", Onset: time.Hour*1, Halflife: (time.Hour*5 + time.Minute*42) }
	go Metabolize(inputDose, caffeine, channel)
	for remnant := range channel {
		inputDose.Remnants = append(inputDose.Remnants, remnant)
	}

	printable, err := json.Marshal(inputDose)
	check(err)

	if *flagWrite == true {
		_, err = f.WriteString(string(printable))
		check(err)
	} else {
		fmt.Println(string(printable))
	}
}
