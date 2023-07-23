package main

import (
	"fmt"
	"log"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Print(r)
	fmt.Fprint(w, "OK")
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	log.Print(r)
	resps, err := callback(r)
	if err != nil {
		if err.Error() == ErrInvalidSignature().Error() || err.Error() == ErrInvalidMessage().Error() {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Print(err)
			return
		} else {
			http.Error(w, "InternalServerError", http.StatusInternalServerError)
			log.Print(err)
			return
		}
	}
	for _, resp := range resps {
		fmt.Fprint(w, resp)
	}
}
