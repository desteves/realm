// Uses standard "sample data" available on all Atlas clusters.
// See https://docs.atlas.mongodb.com/sample-data/
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
	fmt.Printf("Got new graphql client!\n")

	err = client.Connect()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	fmt.Printf("graphql client connected!\n")

	_, err = client.Health()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	fmt.Printf("Passed healthcheck test!\n")


	// Query Sample Data
// TODO

	// Mutate Sample Data
// TODO

	
	fmt.Printf("The End.\n")

}
