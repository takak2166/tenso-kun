package main

import (
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func callback(bot *linebot.Client, r *http.Request) (resp *http.Response, err error) {
	events, err := bot.ParseRequest(r)

	// validate signature
	if err != nil {
		return nil, err
	}

	for _, event := range events {
		userName := "Unknown"
		picURL := "https://imgur.com/ZelRJVU.png"
		if event.Source.UserID != "" {
			userName, picURL, err = getUserInfo(bot, event.Source.UserID)
			if err != nil {
				return nil, err
			}
		}
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				resp, err = sendToWebhook(os.Getenv("WEBHOOK_LINK"), ReqestParam{userName, message.Text, picURL})
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return resp, nil
}
