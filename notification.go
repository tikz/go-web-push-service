package main

import (
	"encoding/json"
	"fmt"

	"github.com/SherClockHolmes/webpush-go"
)

type Notification struct {
	Title   string `json:"title"`
	Body    string `json:"body"`
	IconURL string `json:"icon"`
}

func sendNotification(notification *Notification) {
	payload, err := json.Marshal(notification)
	if err != nil {
		fmt.Println("error JSON marshalling notification:", err)
	}
	for _, sub := range channel.Subscriptions {
		_, err := webpush.SendNotification(payload, sub, channel.Options)
		if err != nil {
			// TODO: Handle error
			fmt.Println("err", err)
		}
	}
}
