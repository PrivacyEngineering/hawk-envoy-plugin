package main

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"log"
)

func send(data []byte, usages []map[string]interface{}) {
	parsed, err := json.Marshal(usages)
	if err != nil {
		log.Printf("Unable to marshall json: %s", err)
		toDlq(data)
		return
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(parsed).
		Post(collectorUrl)

	if err != nil || resp.StatusCode() != 200 {
		log.Printf("Back to dlq for [%d, %s]", resp.StatusCode(), err)
		toDlq(parsed)
	}
}
