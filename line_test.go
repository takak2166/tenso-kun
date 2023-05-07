package main

import (
	"log"
	"os"
	"testing"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func TestGetUserInfo(t *testing.T) {
	testCases := []struct {
		name             string
		userId           string
		expectedUserName string
		expectedPicURL   string
	}{
		{
			name:             "case_me",
			userId:           os.Getenv("TEST_USER_ID"),
			expectedUserName: os.Getenv("TEST_USER_NAME"),
			expectedPicURL:   os.Getenv("TEST_PIC_URL"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bot, err := linebot.New(
				os.Getenv("CHANNEL_SECRET"),
				os.Getenv("CHANNEL_ACCESS_TOKEN"),
			)
			if err != nil {
				t.Error(err)
			}
			userName, picURL, err := getUserInfo(bot, tc.userId)
			if err != nil {
				t.Errorf("getUserInfo() raised error: %v", err)
			}
			if userName != tc.expectedUserName {
				t.Errorf("getUserInfo() expected %v but %v", tc.expectedUserName, userName)
			}
			if picURL != tc.expectedPicURL {
				t.Errorf("getUserInfo() expected %v but %v", tc.expectedPicURL, picURL)
			}
			log.Printf("picURL: %v\n", picURL)
		})
	}
}
