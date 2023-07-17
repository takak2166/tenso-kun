package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/callback", callbackHandler)
	log.Fatal(http.ListenAndServe(":5000", mux))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Print(r)
	fmt.Fprint(w, "OK")
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	log.Print(r)
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_ACCESS_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := callback(bot, r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Print(err)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal(err)
		}
	}
	log.Print(resp)
	fmt.Fprint(w, resp)
}
