package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var (
	errTextInvalidSignature = "InvalidSignatureError"
	errTextInvalidMessage   = "InvalidMessageError"
	defaultUserName         = "Unknown"
	defaultPicURL           = "https://imgur.com/ZelRJVU.png"
)

func ErrInvalidSignature() error {
	return errors.New(errTextInvalidSignature)
}

func ErrInvalidMessage() error {
	return errors.New(errTextInvalidMessage)
}

func callback(r *http.Request) (resp *http.Response, err error) {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_ACCESS_TOKEN"),
	)
	if err != nil {
		return nil, err
	}
	events, err := bot.ParseRequest(r)

	// validate signature
	if err != nil {
		if err.Error() == linebot.ErrInvalidSignature.Error() {
			return nil, ErrInvalidSignature()
		} else {
			return nil, err
		}
	}

	for _, event := range events {
		userName := defaultUserName
		picURL := defaultPicURL
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
			default:
				return nil, ErrInvalidMessage()
			}
		}
	}
	return resp, nil
}
