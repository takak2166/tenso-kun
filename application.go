package main

import (
	"net/http"
	"os"
)

func callback(r *http.Request) ([]*http.Response, error) {
	resps := []*http.Response{}
	messageInfo, err := readEvents(r)
	if err != nil {
		return nil, err
	}
	for _, msgInfo := range messageInfo {
		whReqParam := WebhookReqParam{
			UserName: msgInfo.UserName,
			Text:     msgInfo.Message,
			IconURL:  msgInfo.PicURL,
		}
		resp, err := sendToWebhook(os.Getenv("WEBHOOK_LINK"), whReqParam)
		if err != nil {
			return nil, err
		}
		resps = append(resps, resp)
	}
	return resps, nil
}
