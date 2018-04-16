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

func calculateRemnants(t time.Time, mg int, remnants chan remnant) {
    for inmg := mg/2; inmg >= 1; inmg /= 2 {
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

func marshallRemnant(format string, remrem remnant) (formatted string) {
    if format == "json" {
        line := fmt.Sprintf("     {\"time\": \"%s\", \"remnant\": %d}", remrem.date, remrem.amount)
        if (remrem.amount > 1) {
            line += ","
        }
        return line + "\n"
    } else {
        return fmt.Sprintf("%s,%d\n", remrem.date, remrem.amount)
    }
}

type dose struct {
    name string
    dosage int
    time string
    remnants []remnant
}

func marshallDose(format string, dokidoki dose) (formatted string) {
    if format == "json" {
        return fmt.Sprintf("   {\"name\": \"%s\",\n    \"dosage\": %d\n" +
                   "    \"time\": \"%s\"\n    \"remnants\": [", dokidoki.name, dokidoki.dosage, dokidoki.time)
    } else {
        return fmt.Sprintf("%s,%d", dokidoki.time, dokidoki.dosage)
    }
}

func main() {
    writePtr := flag.Bool("write", false, "write output according to format")
    outputFmt := flag.String("format", "csv", "use format \"csv\" or \"json\"")
    tf := flag.String("time", "now", "time of dosage, e.g. 2018-04-16T17:22:40Z")
    dosaggio := flag.Int("dosage", 100, "caffeine dosage in mg")
    nombre := flag.String("name", "Americano", "name of caffeine product")
    flag.Parse()

    t := time.Now()
    var err error
    if *tf != "now" {
        t, err = time.Parse(time.RFC3339, *tf)
        check(err)
    }
    ts := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02dZ",
            t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

    var f *os.File
    doshi := dose { name: *nombre, dosage: *dosaggio, time: ts, remnants: nil  }
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

    bufSize := nHalves(*dosaggio)
    channel := make(chan remnant, bufSize)
    go calculateRemnants(t, *dosaggio, channel)

    var line string
    for remy := range channel {
        line = marshallRemnant(*outputFmt, remy)

        if *writePtr == true {
            _, err = f.WriteString(line)
            check(err)
        } else {
            fmt.Printf(line)
        }
    }

    if *outputFmt == "json" {
        footer := "   ]\n  }\n]\n"
        if *writePtr == true {
            _, err = f.WriteString(footer)
        } else {
            fmt.Printf(footer)
        }
    }
}
