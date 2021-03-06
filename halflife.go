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

const defaultFilename string = "caff.json"

func LoadDoses() (doses []Dose) {
	jsonFile, err := os.Open(defaultFilename)
	check(err)

	byteValue, err := ioutil.ReadAll(jsonFile)
	check(err)

	jsonFile.Close()

	var ds []Dose
	err = json.Unmarshal(byteValue, &ds)
	check(err)

	return ds
}

func WriteDoses(doses []Dose, toStdout bool) {
        printable, err := json.Marshal(doses)
        check(err)

        if toStdout == false {
                f, err := os.OpenFile(defaultFilename, os.O_CREATE|os.O_WRONLY, 0644)
                check(err)
                _, err = f.WriteString(string(printable))
                check(err)
                f.Close()
        } else {
                fmt.Println(string(printable))
	}
}

const defaultTime string = "2018-04-16T17:22:40Z"

func main() {
	flagWrite := flag.Bool("write", false, "write output to " + defaultFilename)
	flagTime := flag.String("time", defaultTime, "time of dosage, e.g. " + defaultTime)
	flagDosage := flag.Int("dosage", 100, "caffeine dosage in mg")
	flagName := flag.String("name", "Americano", "name of caffeine product")
	flagRead := flag.Bool("read", false, "read json in " + defaultFilename)
	flag.Parse()

	var doses []Dose
	if *flagRead == true {
		doses = LoadDoses()
		fmt.Printf("Read in %d doses from file; input will be appended.\n", len(doses))
	}

        inputDose := Dose{ Name: *flagName, Dosage: *flagDosage, Time: *flagTime, Remnants: nil }
	bufSize := Halves(*flagDosage)
	channel := make(chan Remnant, bufSize)

	caffeine := Metabolizer { Name: "caffeine", Onset: time.Hour*1, Halflife: (time.Hour*5 + time.Minute*42) }
	go Metabolize(inputDose, caffeine, channel)
	for remnant := range channel {
		inputDose.Remnants = append(inputDose.Remnants, remnant)
	}

	doses = append(doses, inputDose)
	WriteDoses(doses, !(*flagWrite))
}
