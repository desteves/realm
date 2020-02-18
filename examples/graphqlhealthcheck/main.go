// Package main a simple graphql instance.
package main

import (
	"fmt"

	g "github.com/desteves/realm/pkg/graphql"
	"github.com/desteves/realm/pkg/options"
	log "github.com/sirupsen/logrus"
)

func main() {

	appid := "graphqlserver-lrnqt" // please don't ddos my poor little app, leaving it open so y'all can test etc.
	auth := "anon-user"
	opts := &options.ClientOptions{
		AppID:         &appid,
		AuthMechanism: &auth}

	client, err := g.NewClient(opts)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	fmt.Printf("Got new GraphQL client!\n")

	err = client.Connect()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	fmt.Printf("GraphQL client connected!\n")

	// this requires a very opinionated graphql configuration ;)
	q, err := client.Health()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	fmt.Printf("Passed healthcheck test, got %+v \n", q)

	// err = client.Disconnect()
	// if err != nil {
	// 	log.Fatalf("%+v", err)
	// }
	fmt.Printf("The End.\n")

}
