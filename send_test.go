package main

import (
	"net/http"
	"os"
	"testing"
)

func TestSendToSlack(t *testing.T) {
	testCases := []struct {
		name        string
		reqestParam ReqestParam
	}{
		{
			name: "case1",
			reqestParam: ReqestParam{
				UserName: "test_user",
				Text:     "test message",
				IconURL:  "https://imgur.com/ZelRJVU.png",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if resp, err := sendToWebhook(os.Getenv("WEBHOOK_LINK"), tc.reqestParam); err != nil {
				t.Errorf("sendToSlack() raised error: %v", err)
			} else if resp.StatusCode != http.StatusOK {
				t.Errorf("sendToSlack() wants ok but %v", resp.StatusCode)
			}
		})
	}
}
