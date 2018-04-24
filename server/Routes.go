package main

import "net/http"

type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
    Route{
        "Index",
        "GET",
        "/",
        Index,
    },
    Route{
        "DoseIndex",
        "GET",
        "/doses",
        DoseIndex,
    },
    Route{
        "DoseShow",
        "GET",
        "/doses/{doseId}",
        DoseShow,
    },
    Route{
        "DoseCreate",
        "POST",
        "/doses",
        DoseCreate,
    },
}

