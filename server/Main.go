package main

import (
    "log"
    "net/http"
    "os"
    "fmt"
)

func determineListenAddress() (string, error) {
  port := os.Getenv("PORT")
  if port == "" {
    return ":8080", fmt.Errorf("$PORT not set")
  }
  return ":" + port, nil
}

func main() {
    router := NewRouter()
    addr, _ := determineListenAddress()

    log.Fatal(http.ListenAndServe(addr, router))
}
