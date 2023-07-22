package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/callback", callbackHandler)
	log.Fatal(http.ListenAndServe(":5000", mux))
}
