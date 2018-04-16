package main

import (
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

func main() {
    writePtr := flag.Bool("write", false, "write csv output to caff.csv")
    flag.Parse()

    t, err := getEntryTime()
    check(err)

    mg, err := getEntryMg()
    check(err)

    if *writePtr == true {
        f, err := os.Create("caffeine.csv")
        check(err)
        defer f.Close()
    }

    fmt.Println("date,close");
    for inmg := mg; inmg >= 1; inmg /= 2 {
        fmt.Printf("%d-%02d-%02dT%02d:%02d:%02dZ,%d\n",
            t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), inmg)

        //The halflife of caffeine is 5.7 hours, aka 5 hours and 42 minutes
        t = t.Add(time.Hour * 5 + time.Minute * 42)
    }
}
