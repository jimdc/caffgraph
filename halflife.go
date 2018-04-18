package main

import (
	"io/ioutil"
	"strconv"
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

// isn't this the standard? (I might have written it custom when it wasn't, 
// then manually changed it back to output the standard time)
func formatTime(t time.Time) (tt string) {
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02dZ",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func CalculateRemnants(t time.Time, mg int, Remnants chan Remnant) {
	for inmg := mg / 2; inmg >= 1; inmg /= 2 {
		Remnants <- Remnant{Date: formatTime(t), Amount: inmg}
		//The halflife of caffeine is 5.7 hours, aka 5 hours and 42 minutes
		t = t.Add(time.Hour*5 + time.Minute*42)
	}
	close(Remnants)
}

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

func marshalDose(dokidoki Dose) (formatted string) {
	return fmt.Sprintf("  {\"name\": \"%s\",\n   \"dosage\": %d,\n"+
		"   \"time\": \"%s\",\n   \"remnants\": [\n", dokidoki.Name, dokidoki.Dosage, dokidoki.Time)
}

var Filename = "caff.json"

func main() {
	writePtr := flag.Bool("write", false, "write output to " + Filename)
	tf := flag.String("time", "now", "time of dosage, e.g. 2018-04-16T17:22:40Z")
	dosaggio := flag.Int("dosage", 100, "caffeine dosage in mg")
	nombre := flag.String("name", "Americano", "name of caffeine product")
	readPtr := flag.Bool("read", false, "read json in " + Filename)
	flag.Parse()

	if *readPtr == true {
		jsonFile, err := os.Open(Filename)
		check(err)
		defer jsonFile.Close()

		byteValue, err := ioutil.ReadAll(jsonFile)
		check(err)

		var ds []Dose
		err = json.Unmarshal(byteValue, &ds)
		check(err)
		//fmt.Printf("%v\n", ds)

		for i := 0; i < len(ds); i++ {
			dis := ds[i]

			fmt.Println("name: " + dis.Name)
			fmt.Println("dosage: " + strconv.Itoa(dis.Dosage))
			fmt.Println("time: " + dis.Time)

			for j := 0; j < len(dis.Remnants); j++ {
				fmt.Println("Remnant time: " + dis.Remnants[j].Date)
				fmt.Println("Remnant Amount: " + strconv.Itoa(dis.Remnants[j].Amount))
			}
		}
	}

	t := time.Now()
	var err error
	if *tf != "now" {
		t, err = time.Parse(time.RFC3339, *tf)
		check(err)
	}
	var f *os.File
	doshi := Dose{Name: *nombre, Dosage: *dosaggio, Time: formatTime(t), Remnants: nil}
	// TODO: check if this is really the beginning of the file and add [ if so. Not otherwise.
	sadoshi := "[\n" + marshalDose(doshi)
	if *writePtr == true {
		f, err = os.OpenFile(Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		check(err)
		defer f.Close()
		// TODO: json.Marshal this dose. After calculating remnants. 
		_, err = f.WriteString(sadoshi)
		check(err)
	} else {
		fmt.Println(sadoshi)
	}

	bufSize := Halves(*dosaggio)
	channel := make(chan Remnant, bufSize)
	go CalculateRemnants(t, *dosaggio, channel)

	var line string;
	for remy := range channel {

		remypretty, _ := json.MarshalIndent(remy, "", "    ")
		line = string(remypretty)

		bufSize--
		if bufSize > 1 {
			line += ","
		}

		if *writePtr == true {
			_, err = f.WriteString(line+"\n")
			check(err)
		} else {
			fmt.Println(line)
		}
	}

	if footer := "   ]\n  }\n]\n";
	   *writePtr == true {
		_, err = f.WriteString(footer)
		check(err)
	} else {
		fmt.Printf(footer)
	}
}
