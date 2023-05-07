package main

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type ReqestParam struct {
	UserName string `json:"username"`
	Text     string `json:"text"`
	IconURL  string `json:"icon_url"`
}

func sendToWebhook(webhookLink string, reqParam ReqestParam) (*http.Response, error) {
	p, err := json.Marshal(reqParam)
	if err != nil {
		return nil, err
	}
	resp, err := http.PostForm(webhookLink, url.Values{"payload": {string(p)}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp, nil
}
