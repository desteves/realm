// Package main a simple graphql instance.
package main

import (
	"fmt"
	"log"

	r "github.com/desteves/realm/pkg/auth"
	"github.com/desteves/realm/pkg/options"
)

func main() {

	appid := "graphqlserver-lrnqt" // please don't ddos my poor little app, leaving it open so y'all can test etc.
	auth := "anon-user"
	opts := &options.ClientOptions{
		AppID:         &appid,
		AuthMechanism: &auth}
	client, err := r.NewClient(opts)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	fmt.Printf("Got new client!\n")

	err = client.Connect()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	fmt.Printf("Client connected to Realm!\n")

	err = client.Ping()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	fmt.Printf("Passed webhook ping test!\n")

	err = client.Disconnect()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	fmt.Printf("The End.\n")
}
