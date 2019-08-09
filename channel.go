package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"

	"github.com/SherClockHolmes/webpush-go"
)

func createChannel() {
	log.Println("Generating VAPID key pair...")
	file, err := os.Create(channelFile)
	defer file.Close()
	if err != nil {
		log.Fatalln("error creating channel file:", err)
	} else {
		privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
		if err != nil {
			log.Fatalln("error generating VAPID keys:", err)
		}

		newChannel := Channel{Options: &webpush.Options{
			VAPIDPublicKey:  publicKey,
			VAPIDPrivateKey: privateKey,
			TTL:             30,
		}}

		encoder := gob.NewEncoder(file)
		encoder.Encode(&newChannel)
	}
	fmt.Println("New channel saved in", channelFile)
}

func saveChannel() {
	file, err := os.Create(channelFile)
	defer file.Close()
	if err != nil {
		log.Fatalln("error creating channel file:", err)
	} else {
		encoder := gob.NewEncoder(file)
		encoder.Encode(&channel)
	}
}

func loadChannel() {
	if _, err := os.Stat(channelFile); os.IsNotExist(err) {
		log.Println(channelFile, "not found.")
		createChannel()
	}
	file, err := os.Open(channelFile)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(&channel)
	}
	file.Close()
}
