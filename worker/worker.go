package main

import "fmt"
import "log"
import "net/http"

func main() {
    fmt.Println("hello worker")

    http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "pong")
    })

    log.Fatal(http.ListenAndServe(":5060", nil))
}
