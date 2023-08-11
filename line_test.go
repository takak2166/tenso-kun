package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var callbackTestRequestBody = `{
    "events":[
		{
			"replyToken":"00000000000000000000000000000000",
			"type":"message",
			"timestamp":1570728501562,
			"source":{
				"type":"user",
				"userId":"Udeadbeefdeadbeefdeadbeefdeadbeef"
			},
			"message":{
				"id":"100001",
				"type":"text",
				"text":"Hello, world"
			}
		},
		{
			"replyToken":"ffffffffffffffffffffffffffffffff",
			"type":"message",
			"timestamp":1570728501562,
			"source":{
				"type":"user",
				"userId":"Udeadbeefdeadbeefdeadbeefdeadbeef"
			},
			"message":{
				"id":"100002",
				"type":"sticker",
				"packageId":"1",
				"stickerId":"1"
			}
		}
    ]
}`

func genLineSignature(secretKey string, requestBody string) string {
	hash := hmac.New(sha256.New, []byte(secretKey))
	hash.Write([]byte(requestBody))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func TestReadEvents(t *testing.T) {
	testCases := []struct {
		name             string
		valid            bool
		requestBody      []byte
		xLineSignature   string
		expectedMsgInfo1 MsgInfo
		expectedMsgInfo2 MsgInfo
	}{
		{
			name:             "valid case 1",
			valid:            true,
			requestBody:      []byte(callbackTestRequestBody),
			xLineSignature:   genLineSignature(os.Getenv("CHANNEL_SECRET"), callbackTestRequestBody),
			expectedMsgInfo1: MsgInfo{"Unknown", "Hello, world", "https://imgur.com/ZelRJVU.png"},
			expectedMsgInfo2: MsgInfo{"Unknown", "Unknown Message", "https://imgur.com/ZelRJVU.png"},
		},
		{
			name:           "valid case 2",
			valid:          true,
			requestBody:    []byte(`{"events":[]}`),
			xLineSignature: genLineSignature(os.Getenv("CHANNEL_SECRET"), `{"events":[]}`),
		},
		{
			name:           "invalid case",
			valid:          false,
			requestBody:    []byte(callbackTestRequestBody),
			xLineSignature: genLineSignature("hoge", callbackTestRequestBody),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/callback", bytes.NewReader(tc.requestBody))
			req.Header.Set("X-Line-Signature", tc.xLineSignature)
			msgInfo, err := readEvents(req)
			if err != nil {
				if !tc.valid {
					return
				}
				t.Errorf("readEvents() raised error: %v", err)
				return
			}
			t.Log(msgInfo)
		})
	}
}

func TestGetUserInfo(t *testing.T) {
	testCases := []struct {
		name             string
		userId           string
		valid            bool
		expectedUserName string
		expectedPicURL   string
	}{
		{
			name:             "case_me",
			userId:           os.Getenv("TEST_USER_ID"),
			valid:            true,
			expectedUserName: os.Getenv("TEST_USER_NAME"),
			expectedPicURL:   os.Getenv("TEST_PIC_URL"),
		},
		{
			name:             "not found",
			userId:           "Udeadbeefdeadbeefdeadbeefdeadbeef",
			valid:            true,
			expectedUserName: defaultUserName,
			expectedPicURL:   defaultPicURL,
		},
		{
			name:             "no userId",
			userId:           "",
			valid:            true,
			expectedUserName: defaultUserName,
			expectedPicURL:   defaultPicURL,
		},
		{
			name:             "invalid userId",
			userId:           "!@#$%^&*()",
			valid:            false,
			expectedUserName: "",
			expectedPicURL:   "",
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
				if !tc.valid {
					return
				}
				t.Errorf("getUserInfo() raised error: %v", err)
				return
			}
			if userName != tc.expectedUserName {
				t.Errorf("getUserInfo() expected %v but %v", tc.expectedUserName, userName)
				return
			}
			if picURL != tc.expectedPicURL {
				t.Errorf("getUserInfo() expected %v but %v", tc.expectedPicURL, picURL)
				return
			}
		})
	}
}
