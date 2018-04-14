package main

import (
    "fmt"
    "bufio"
    "os"
    "strconv"
    "strings"
    "time"
)

func main() {
    const halflife = 5.7
    reader := bufio.NewReader(os.Stdin)
    t := time.Now()

    fmt.Print("Enter mg: ")
    str, _ := reader.ReadString('\n')
    mg, err := strconv.ParseInt(strings.TrimSpace(str), 10, 64)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    for inmg := mg; inmg >= 1; inmg /= 2 {
        fmt.Printf("%d-%02d-%02dT%02d:%02d %dmg\n",
            t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), inmg)

        //The halflife of caffeine is 5.7 hours, aka 5 hours and 42 minutes
        t = t.Add(time.Hour * 5 + time.Minute * 42)
    }
}
