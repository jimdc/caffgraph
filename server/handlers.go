package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "html/template"
    "io"
    "io/ioutil"
    "github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Welcome!")
}

func DoseIndex(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(doses); err != nil {
        panic(err)
    }
}

func DoseShow(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    doseId := vars["doseId"]
    fmt.Fprintln(w, "Dose show:", doseId)
}

func DoseCreate(w http.ResponseWriter, r *http.Request) {
    var dose Dose
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }
    if err := json.Unmarshal(body, &dose); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    }

    d := RepoCreateDose(dose)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(d); err != nil {
        panic(err)
    }
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
    p, err := loadPage("Title") //so it looks for Title.txt
    if err != nil {
        fmt.Println(err)
    }
    t, err := template.ParseFiles("test.html")
    if err != nil {
        fmt.Println(err)
    }
    t.Execute(w, p)
}
