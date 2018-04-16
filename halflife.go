package main

import (
    "math"
    "flag"
    "fmt"
    "os"
    "time"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

type remnant struct {
    date string
    amount int
}

// math.Round doesn't round up, so compensate when allocating
func nHalves(mg int) (timesToDivide int) {
   fmg := float64(mg)
   nDivides := math.Round(math.Log2(fmg))
   return int(nDivides)
}

func calculateHl(t time.Time, mg int, remnants chan remnant) {
    for inmg := mg; inmg >= 1; inmg /= 2 {
        formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02dZ",
            t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

        remnants <- remnant { date: formatted, amount: inmg }

        //The halflife of caffeine is 5.7 hours, aka 5 hours and 42 minutes
        t = t.Add(time.Hour * 5 + time.Minute * 42)
    }
    close(remnants)
}

func header(format string) (headerFormatted string) {
    if format == "json" {
        return "[";
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

func main() {
    writePtr := flag.Bool("write", false, "write output to caffeine.csv")
    outputFmt := flag.String("format", "csv", "use format \"csv\" or \"json\"")
    ts := flag.String("time", "now", "time of dosage, e.g. 2018-04-16T17:22:40Z")
    dose := flag.Int("dose", 100, "caffeine dosage in mg")
    flag.Parse()

    t := time.Now()
    var err error
    if *ts != "now" {
        t, err = time.Parse(time.RFC3339, *ts)
        check(err)
    }

    var f *os.File
    if *writePtr == true {
        f = opendata(*outputFmt)
        defer f.Close()
    } else {
        fmt.Println(header(*outputFmt))
    }

    bufSize := nHalves(*dose) + 1
    channel := make(chan remnant, bufSize)
    go calculateHl(t, *dose, channel)

    var line string
    for remy := range channel {

        if *outputFmt == "json" {
            line = fmt.Sprintf("  {\"time\": \"%s\", \"remnant\": %d}\n", remy.date, remy.amount)
        } else {
            line = fmt.Sprintf("%s,%d\n", remy.date, remy.amount)
        }

        if *writePtr == true {
            _, err = f.WriteString(line)
            check(err)
        } else {
            fmt.Printf(line)
        }
    }

    if *outputFmt == "json" {
        footer := "]\n"
        if *writePtr == true {
            _, err = f.WriteString(footer)
        } else {
            fmt.Printf(footer)
        }
    }
}
