package main

import (
    "log"
    "net/http"
    "os"
    "io/ioutil"
    "fmt"
)

func determineListenAddress() (string, error) {
  port := os.Getenv("PORT")
  if port == "" {
    return ":8080", fmt.Errorf("$PORT not set")
  }
  return ":" + port, nil
}

/*
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
func testHandler(w http.ResponseWriter, r *http.Request, title string)
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }
        fn(w, r, m[2])
    }
}
*/

type Page struct {
    Title string
    Body  []byte
}

func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

func main() {
    router := NewRouter()
    addr, _ := determineListenAddress()

    log.Fatal(http.ListenAndServe(addr, router))
}
