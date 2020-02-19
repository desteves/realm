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

type jsondict map[string]interface{}

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

func runQuery(client *g.Client, namespace string, query interface{}, variables map[string]interface{}) {

	var response g.Response
	fmt.Printf("\n%v Query ", namespace)
	err := client.Query(context.TODO(), query, variables, &response)
	if err != nil {
		fmt.Printf("! err %+v \n", err.Error())
		return
	}
	fmt.Printf("Response %+v \n", response)
}

func runMutation(client *g.Client, namespace string, mutation string, variables map[string]interface{}) {

	var response g.Response
	fmt.Printf("\n%v Mutation ", namespace)
	err := client.Mutate(context.TODO(), mutation, variables, &response)
	if err != nil {
		fmt.Printf("! err %+v \n", err.Error())
		return
	}
	fmt.Printf("Response %+v \n", response)
}

func runSampleQueries(client *g.Client) {

	// findOne
	var qOne struct {
		AirBnBListingAndReview struct {
			ID   string `graphql:"_id"`
			Name string `graphql:"name"`
			URI  string `graphql:"listing_url"`
		} `graphql:"listingsAndReviews"`
	}
	runQuery(client, "sample_airbnb.listingsAndReviews", qOne, nil)

	// findOne with filter
	var qTwo struct {
		AirBnBListingAndReview struct {
			ID   string `graphql:"_id"`
			Name string `graphql:"name"`
			URI  string `graphql:"listing_url"`
		} `graphql:"listingsAndReviews( query: { _id: \"10009999\" } )"`
	}
	runQuery(client, "sample_airbnb.listingsAndReviews", qTwo, nil)

	// find many with limit
	var qThree struct {
		AirBnBListingAndReviews []struct { // don't forget slice, else error: "slice doesn't exist in any of 1 places to unmarshal"
			ID   string `graphql:"_id"`
			Name string `graphql:"name"`
			URI  string `graphql:"listing_url"`
		} `graphql:"listingsAndReviewss( limit: 3 )"`
	}
	runQuery(client, "sample_airbnb.listingsAndReviews", qThree, nil)

	// find many with sort
	var qFour struct {
		Accounts []struct {
			ID        string `graphql:"_id"`
			AccountID string `graphql:"account_id"`
		} `graphql:"accountss(sortBy: ACCOUNT_ID_ASC, limit: 5)"`
	}
	runQuery(client, "sample_analytics.accounts", qFour, nil)

}

func runSampleMutations(client *g.Client) {

	// insertOne<collection> and return _id
	mOne := `mutation ($customer: CustomerInsertInput!) {
		insertOneCustomer(data: $customer) {
			_id
		}	}`

	vOne := g.Variable{
		"customer": g.Variable{
			"active":   false,
			"address":  "123 4th st apt 5",
			"name":     "diana",
			"username": "d",
		},
	}
	runMutation(client, "sample_analytics.customers", mOne, vOne)

	// insertMany<collection>s and return _id's
	mTwo := `mutation ($customers: [CustomerInsertInput!]!) {
		insertManyCustomers(data: $customers) {
			insertedIds
		}
	}`

	vTwo := g.Variable{
		"customers": []jsondict{
			jsondict{
				"active":   false,
				"address":  "123 4th st apt 5",
				"name":     "diana",
				"username": "d",
			},
			jsondict{
				"active":   false,
				"address":  "33 sheridan st",
				"name":     "griffin",
				"username": "the_cat",
			},
		},
	}
	runMutation(client, "sample_analytics.customers", mTwo, vTwo)

	// deleteMany<collection>s and return _id's
	mThree := `mutation ($query: TransactionQueryInput) {
		deleteManyTransactions(query: $query) {
			deletedCount
		}
	}`
	vThree := g.Variable{
		"query": jsondict{
			"account_id": 278603, //  996263, 443178, 716662, 996263
		},
	}
	runMutation(client, "sample_analytics.transactions", mThree, vThree)

	// deleteOne<collection>  and return _id
	mFour := `mutation ($query: TransactionQueryInput!) {
		deleteOneTransaction(query: $query) {
			_id
		}
	}`
	vFour := g.Variable{
		"query": jsondict{
			"transaction_count": 40,
		},
	}
	runMutation(client, "sample_training.tweets", mFour, vFour)

	// replaceOne<collection> and return _id
	mFive := `mutation ($query: MovieQueryInput, $data: MovieInsertInput!) {
		replaceOneMovie(query: $query, data: $data) {
			_id
		}
	}`
	vFive := g.Variable{
		"query": jsondict{
			"year": 1940,
		},
		"data": jsondict{
			"tomatoes": jsondict{
				"consensus": "",
				"website":   "stitch.mongodb.com",
			},
			"year":  2020,
			"type":  "horror",
			"title": "learning graphql",
		},
	}
	runMutation(client, "sample_mflix.movies", mFive, vFive)

	// updateMany<collection>s and return matched & modified count
	mSix := `mutation ($query: CommentQueryInput, $set: CommentUpdateInput!) {
		updateManyComments(query: $query, set: $set) {
			matchedCount
			modifiedCount
		}
	}`
	vSix := g.Variable{
		"query": jsondict{
			"name": "Jon Snow",
		},
		"set": jsondict{
			"name": "Aegon Targaryen",
		},
	}
	runMutation(client, "sample_mflix.comments", mSix, vSix)

	// updateOne<collection> and return _id
	mSeven := `mutation ($query: UserQueryInput, $set: UserUpdateInput!) {
		updateOneUser(query: $query, set: $set) {
			_id
			name
			email
		}
	}`
	vSeven := g.Variable{
		"query": jsondict{
			"name": "Jon Snow",
		},
		"set": jsondict{
			"name": "Aegon Targaryen",
		},
	}

	runMutation(client, "sample_mflix.users", mSeven, vSeven)

	// upsertOne<collection> and return _id
	mEight := `mutation ($q: TheaterQueryInput, $d: TheaterInsertInput!) {
		upsertOneTheater(query: $q, data: $d) {
			_id
			location {
				address {
					city
					state
					zipcode
					street1
					street2
				}
			}
		}
	}`
	vEight := g.Variable{
		"q": jsondict{
			"theaterId": 99999,
		},
		"d": jsondict{
			"location": jsondict{
				"address": jsondict{
					"zipcode": "78701",
					"city":    "Austin",
					"state":   "TX",
					"street1": "1800 Congress Ave",
				},
			},
		},
	}
	runMutation(client, "sample_supplies.sales", mEight, vEight)

}

// GraphiQL syntax for mutations
//
// mutation ($customer: CustomerInsertInput!) {
// 	insertOneCustomer(data: $customer) {
// 		_id
// 	}
// }
// {
//   "customer": {
//     "active": false,
//     "address": "123 4th st apt 5",
//     "name": "diana",
//     "username": "d"
//   }
// }
