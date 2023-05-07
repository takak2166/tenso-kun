package main

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func getUserInfo(bot *linebot.Client, userId string) (string, string, error) {
	userProf, err := bot.GetProfile(userId).Do()
	if err != nil {
		return "", "", err
	}
	return userProf.DisplayName, userProf.PictureURL, nil
}
