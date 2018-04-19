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

const Filename string = "caff.json"

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
        doshi := Dose{Name: *flagName, Dosage: *flagDosage, Time: ts, Remnants: nil}

	//TODO: integrate with write, to not have duplicate arrays and stuff
	if *flagRead == true {
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

	var f *os.File
	if *flagWrite == true {
		f, err = os.OpenFile(Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		check(err)
		defer f.Close()
	}

	bufSize := Halves(*flagDosage)
	channel := make(chan Remnant, bufSize)
	go CalculateRemnants(doshi, channel)
	for residuo := range channel {
		doshi.Remnants = append(doshi.Remnants, residuo)
	}

	printable, err := json.Marshal(doshi)
	check(err)

	if *flagWrite == true {
		_, err = f.WriteString(string(printable))
		check(err)
	} else {
		fmt.Println(string(printable))
	}
}
