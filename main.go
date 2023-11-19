package main

import (
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, chat!"))
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}
