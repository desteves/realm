// Uses standard "sample data" available to all Atlas clusters.
// See https://docs.atlas.mongodb.com/sample-data/
package main

import (
	"context"
	"fmt"
	"log"

	g "github.com/desteves/realm/pkg/graphql"
	"github.com/desteves/realm/pkg/options"
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

	var response g.Response
	err = client.Health(&response)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	fmt.Printf("Passed healthcheck test!\n")

	runSampleQueries(client)
	runSampleMutations(client)

	fmt.Printf("The End.\n")

}

func runQuery(client *g.Client, namespace string, query interface{}, variables interface{}) {

	var response g.Response
	fmt.Printf("\n%v Query", namespace)
	err := client.Query(context.TODO(), query, nil, &response)
	if err != nil {
		return
	}
	fmt.Printf("Response %+v\n", response)
}

func runSampleQueries(client *g.Client) {

	// findOne
	var qOne struct {
		AirBnBListingAndReviews struct {
			ID   string `graphql:"_id"`
			Name string `graphql:"name"`
			URI  string `graphql:"listing_url"`
		} `graphql:"listingsAndReviews"`
	}
	runQuery(client, "sample_airbnb.listingsAndReviews", qOne, nil)

	// findOne with filter
	var qTwo struct {
		AirBnBListingAndReviews struct {
			ID   string `graphql:"_id"`
			Name string `graphql:"name"`
			URI  string `graphql:"listing_url"`
		} `graphql:"listingsAndReviews( name: 5 )"`
	}
	runQuery(client, "sample_airbnb.listingsAndReviews", qTwo, nil)

	// find with limit
	var qThree struct {
		AirBnBListingAndReviewss []struct { // don't forget slice, else error: "slice doesn't exist in any of 1 places to unmarshal"
			ID   string `graphql:"_id"`
			Name string `graphql:"name"`
			URI  string `graphql:"listing_url"`
		} `graphql:"listingsAndReviewss( limit: 5 )"`
	}
	runQuery(client, "sample_airbnb.listingsAndReviews", qThree, nil)

	// filter
	var qFour struct {
		Accounts struct {
			ID        string `graphql:"_id"`
			AccountID string `graphql:"account_id"`
			Limit     string `graphql:"limit"`
		} `graphql:"accounts(query: {limit: 10000}) "`
	}
	runQuery(client, "sample_analytics.accounts", qFour, nil)

}

func runSampleMutations(client *g.Client) {
}

// 	//////////////////////////////////////////////////////////////////////////////////
// 	// Mutate Sample Data from various Databases showcasing a number of filters
// 	//////////////////////////////////////////////////////////////////////////////////
// 	fmt.Printf("Insert a Single Document from sample_analytics.customers\n")
// 	var m struct {
// 		InsertOneCustomer struct {
// 			active   graphql.Boolean
// 			address  graphql.String
// 			email    graphql.String
// 			name     graphql.String
// 			username graphql.String
// 		} `graphql:"insertOneCustomer(data: {
// 			title: "Little Women"
// 			director: "Greta Gerwig"
// 			year: 2019
// 			runtime: 135
// 		})"`
// 	}
// 	variables := map[string]interface{}{
// 		"ep": starwars.Episode("JEDI"),
// 		"review": starwars.ReviewInput{
// 			Stars:      graphql.Int(5),
// 			Commentary: graphql.String("This is a great movie!"),
// 		},
// 	}
// 	err := client.Mutate(context.Background(), &m, variables)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("Created a %v star review: %v \n\n\n", m.CreateReview.Stars, m.CreateReview.Commentary)
// }
