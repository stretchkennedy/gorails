package main

import (
	"log"
	"net/http"
	"os"

	"gopkg.in/igm/sockjs-go.v2/sockjs"
)

func serveSock() {
	var authHandler *AuthHandler

	sockHandler := sockjs.NewHandler(
		"/echo",
		sockjs.DefaultOptions,
		func(session sockjs.Session){
			handleClient(session, authHandler.UserId)
		})

	authHandler = NewAuthHandler(os.Getenv("SECRET"), sockHandler)

	log.Print("attempting to listen...")
	log.Fatal(http.ListenAndServe(os.Getenv("ADDRESS"), authHandler))
}

func handleClient(session sockjs.Session, userId int) {
	log.Printf("listening on behalf of user %d...", userId)

	// recv into channel
	for {
		msg, err := session.Recv()
		if err != nil {
			break // stop listening on error
		}
		session.Send(msg)
	}
}
