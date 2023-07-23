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

type MsgInfo struct {
	UserName string
	Message  string
	PicURL   string
}

func ErrInvalidSignature() error {
	return errors.New(errTextInvalidSignature)
}

func ErrInvalidMessage() error {
	return errors.New(errTextInvalidMessage)
}

func readEvents(r *http.Request) ([]MsgInfo, error) {
	messageInfo := []MsgInfo{}
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
				messageInfo = append(messageInfo, MsgInfo{userName, message.Text, picURL})
			default:
				messageInfo = append(messageInfo, MsgInfo{userName, "Unknown Message", picURL})
			}
		}
	}
	return messageInfo, nil
}

func getUserInfo(bot *linebot.Client, userId string) (string, string, error) {
	userProf, err := bot.GetProfile(userId).Do()
	if err != nil {
		return "", "", err
	}
	return userProf.DisplayName, userProf.PictureURL, nil
}
