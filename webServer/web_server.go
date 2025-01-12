package main

import (
    "fmt"
    //"io"
    "net/http"
    //"os"
)

func main() {
    fs := http.FileServer(http.Dir("pub"))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/" {
            http.ServeFile(w, r, "pub/index.html")
        } else {
            fs.ServeHTTP(w, r)
        }
    })

    fmt.Println("Server listening on port 4000")
    http.ListenAndServe(":4000", nil)
}
