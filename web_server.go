package main

import (
    "net/http"
)

func main() {
    fs := http.FileServer(http.Dir("pub"))
    http.Handle("/", fs)
    http.ListenAndServe(":4000", nil)
}

