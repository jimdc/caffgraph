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

func CalculateRemnants(t time.Time, mg int, Remnants chan Remnant) {
	for inmg := mg / 2; inmg >= 1; inmg /= 2 {
		formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02dZ",
			t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

		Remnants <- Remnant{Date: formatted, Amount: inmg}

		//The halflife of caffeine is 5.7 hours, aka 5 hours and 42 minutes
		t = t.Add(time.Hour*5 + time.Minute*42)
	}
	close(Remnants)
}

func header(format string) (headerFormatted string) {
	if format == "json" {
		return "["
	} else {
		return "date,close"
	}
}

func opendata(format string) (file *os.File) {
	var filename string
	if format == "json" {
		fmt.Println("using json format")
		filename = "caff.json"
	} else {
		fmt.Println("using csv format")
		filename = "caff.csv"
	}

	f, err := os.Create(filename)
	check(err)
	_, err = f.WriteString(header(format))
	check(err)

	return f
}

func marshallRemnant(format string, remrem Remnant) (formatted string) {
	if format == "json" {
		line := fmt.Sprintf("    {\"time\": \"%s\", \"remnant\": %d}", remrem.Date, remrem.Amount)
		if remrem.Amount > 1 {
			line += ","
		}
		return line
	} else {
		return fmt.Sprintf("%s,%d", remrem.Date, remrem.Amount)
	}
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

func marshallDose(format string, dokidoki Dose) (formatted string) {
	if format == "json" {
		return fmt.Sprintf("  {\"name\": \"%s\",\n   \"dosage\": %d,\n"+
			"   \"time\": \"%s\",\n   \"remnants\": [", dokidoki.Name, dokidoki.Dosage, dokidoki.Time)
	} else {
		return fmt.Sprintf("%s,%d", dokidoki.Time, dokidoki.Dosage)
	}
}

func main() {
	writePtr := flag.Bool("write", false, "write output according to format")
	outputFmt := flag.String("format", "csv", "use format \"csv\" or \"json\"")
	tf := flag.String("time", "now", "time of dosage, e.g. 2018-04-16T17:22:40Z")
	dosaggio := flag.Int("dosage", 100, "caffeine dosage in mg")
	nombre := flag.String("name", "Americano", "name of caffeine product")
	readPtr := flag.Bool("read", false, "read json in caff.json")
	flag.Parse()

	if *readPtr == true {
		jsonFile, err := os.Open("caff.json")
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
	ts := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02dZ",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	var f *os.File
	doshi := Dose{Name: *nombre, Dosage: *dosaggio, Time: ts, Remnants: nil}
	sadoshi := marshallDose(*outputFmt, doshi)
	if *writePtr == true {
		f = opendata(*outputFmt)
		defer f.Close()
		_, err = f.WriteString(sadoshi)
		check(err)
	} else {
		fmt.Println(header(*outputFmt))
		fmt.Println(sadoshi)
	}

	bufSize := Halves(*dosaggio)
	channel := make(chan Remnant, bufSize)
	go CalculateRemnants(t, *dosaggio, channel)

	var line string
	for remy := range channel {
		//It works, but annoyingly uses newlines and no space before {
		//remypretty, _ := json.MarshalIndent(remy, "", "    ")
		//line = string(remypretty)
		line = marshallRemnant(*outputFmt, remy)

		if *writePtr == true {
			_, err = f.WriteString(line)
			check(err)
		} else {
			fmt.Println(line)
		}
	}

	if footer := "   ]\n  }\n]\n";
	   *outputFmt == "json" {
		if *writePtr == true {
			_, err = f.WriteString(footer)
			check(err)
		} else {
			fmt.Printf(footer)
		}
	}
}
