package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func main() {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_ACCESS_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "OK")
	})
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		events, err := bot.ParseRequest(r)

		// validate signature
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				http.Error(w, err.Error(), http.StatusBadRequest)
				log.Print(err)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Fatal(err)
			}
		}

		for _, event := range events {
			userName := "Unknown"
			picURL := "https://imgur.com/ZelRJVU.png"
			if event.Source.UserID != "" {
				userName, picURL, err = getUserInfo(bot, event.Source.UserID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					log.Fatal(err)
				}
			}
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					resp, err := sendToWebhook(os.Getenv("WEBHOOK_LINK"), ReqestParam{userName, message.Text, picURL})
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						log.Fatal(err)
					}
					log.Print(resp)
					fmt.Fprint(w, resp)
				}
			}
		}

	})
	log.Fatal(http.ListenAndServe(":5000", mux))
}
