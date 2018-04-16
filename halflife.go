package main

import (
    "math"
    "flag"
    "fmt"
    "bufio"
    "os"
    "strconv"
    "strings"
    "time"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func getEntryTime() (t time.Time, err error) {
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("Enter time: ")
    str, _ := reader.ReadString('\n')
    str = strings.TrimSpace(str)

    if (str != "") {
        t1, e := time.Parse(time.RFC3339, str)
        if e != nil {
            return time.Time{}, e
        }

        return t1, nil
    }

    fmt.Println("Using default current time.")
    return time.Now(), nil
}

func getEntryMg() (mg int64, err error) {
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("Enter mg: ")
    str, _ := reader.ReadString('\n')
    str = strings.TrimSpace(str)

    if (str != "") {
        userMg, e := strconv.ParseInt(strings.TrimSpace(str), 10, 64)
        if e != nil {
            return -1, e
        }

        return userMg, nil
    }

    fmt.Println("Using default 100mg.")
    return 100, nil
}

type remnant struct {
    date string
    amount int64
}

// math.Round doesn't round up, so compensate when allocating
func nHalves(mg int64) (timesToDivide int64) {
   fmg := float64(mg)
   nDivides := math.Round(math.Log2(fmg))
   return int64(nDivides)
}

func calculateHl(t time.Time, mg int64, remnants chan remnant) {
    for inmg := mg; inmg >= 1; inmg /= 2 {
        formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02dZ",
            t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

        remnants <- remnant { date: formatted, amount: inmg }

        //The halflife of caffeine is 5.7 hours, aka 5 hours and 42 minutes
        t = t.Add(time.Hour * 5 + time.Minute * 42)
    }
    close(remnants)
}

func main() {
    writePtr := flag.Bool("write", false, "write csv output to caffeine.csv")
    flag.Parse()
    var f *os.File

    t, err := getEntryTime()
    check(err)

    mg, err := getEntryMg()
    check(err)

    if *writePtr == true {
        f, err = os.Create("caffeine.csv")
        check(err)
        _, err2 := f.WriteString("date,close\n")
        check(err2)
        defer f.Close()
    } else {
        fmt.Println("date,close")
    }

    bufSize := nHalves(mg) + 1
    c := make(chan remnant, bufSize)
    go calculateHl(t, mg, c)
    for remy := range c {
        csv := fmt.Sprintf("%s,%d\n", remy.date, remy.amount)
        if *writePtr == true {
            _, err3 := f.WriteString(csv)
            check(err3)
        } else {
            fmt.Printf(csv)
        }
    }
}
