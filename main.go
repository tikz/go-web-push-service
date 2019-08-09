package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"os"

	webpush "github.com/SherClockHolmes/webpush-go"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

const channelFile = "./channel.gob"

var channel Channel

type Channel struct {
	Subscriptions []*webpush.Subscription
	Options       *webpush.Options
}

func addSubscription(sub *webpush.Subscription) {
	for _, s := range channel.Subscriptions {
		if s.Endpoint == sub.Endpoint {
			return
		}
	}
	channel.Subscriptions = append(channel.Subscriptions, sub)
	saveChannel()
}

func subscribe(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var sub webpush.Subscription
	err := decoder.Decode(&sub)
	if err != nil {
		log.Fatalln("error decoding subscription JSON:", err)
	}

	addSubscription(&sub)
}

func send(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	notification := Notification{
		Title:   r.Form.Get("title"),
		Body:    r.Form.Get("body"),
		IconURL: r.Form.Get("icon"),
	}
	sendNotification(&notification)
}

func publicKey(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(channel.Options.VAPIDPublicKey))
}

func main() {
	loadChannel()

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(os.Getenv("WEB_PUSH_SERVICE_DOMAIN")),
		Cache:      autocert.DirCache("cache/"),
	}
	server := &http.Server{
		Addr: ":https",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
			NextProtos:     []string{acme.ALPNProto},
		},
	}
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/subscribe", subscribe)
	http.HandleFunc("/publicKey", publicKey)
	http.HandleFunc("/send", send)

	go func() {
		h := certManager.HTTPHandler(nil)
		log.Fatal(http.ListenAndServe(":http", h))
	}()

	log.Fatal(server.ListenAndServeTLS("", ""))
}
