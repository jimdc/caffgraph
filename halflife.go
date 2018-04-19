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

// TODO: right now just using a strict average. Better to have a range? (upper, lower)
type Alimentarius struct {
        Name     string        `json:"name"`
        Onset    time.Duration `json:"onset"`
        Halflife time.Duration `json:"halflife"`
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

func CalculateRemnants(damage Dose, Remnants chan Remnant) {

	kafeiyin := Alimentarius { Name: "caffeine", Onset: time.Hour*1, Halflife: (time.Hour*5 + time.Minute*42) }
	// niguding := Alimentarius { Name: "nicotine", Onset: time.Second*5, Halflife: time.Hour*2 }
	// naipusheng := Alimentarius { Name: "naproxen", Onset: time.Minute*30, Halflife: time.Hour*13 }

	t, err := time.Parse(time.RFC3339, damage.Time)
	check(err)
	ts := damage.Time

	for remaining := damage.Dosage / 2; remaining >= 1; remaining /= 2 {
		Remnants <- Remnant{Date: ts, Amount: remaining}

		t = t.Add(kafeiyin.Halflife)
		ts = t.Format(time.RFC3339)
	}

	close(Remnants)
}

func marshalDose(dokidoki Dose) (formatted string) {
	return fmt.Sprintf("  {\"name\": \"%s\",\n   \"dosage\": %d,\n"+
		"   \"time\": \"%s\",\n   \"remnants\": [\n", dokidoki.Name, dokidoki.Dosage, dokidoki.Time)
}

const Filename string = "caff.json"

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

	// TODO: see if it's possible to hide some value in struct from json parser? All this conversion seems ineff
	t := time.Now()
	var err error
	if *tf != "now" {
		t, err = time.Parse(time.RFC3339, *tf)
		check(err)
	}
	ts := t.Format(time.RFC3339)
	doshi := Dose{Name: *nombre, Dosage: *dosaggio, Time: ts, Remnants: nil}
	// TODO: check if this is really the beginning of the file and add [ if so. Not otherwise.
	sadoshi := "[\n" + marshalDose(doshi)
	var f *os.File
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
	go CalculateRemnants(doshi, channel)

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
