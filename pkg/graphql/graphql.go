package graphql

import (
	"context"
	"fmt"
	"log"

	"github.com/desteves/realm/pkg/auth"
	"github.com/desteves/realm/pkg/options"
	gql "github.com/shurcooL/graphql"
)

//HealthCheck checks graphql returns schema
type HealthCheck struct {
	ID          string
	Description string
	Status      string
	Endpoint    string
}

// Client is a Realm GraphQL Client with authentication to a Realm Application
type Client struct {
	client    *gql.Client
	realmAuth *auth.Client
	uri       *string
}

// NewClient creates a new Client
func NewClient(opts *options.ClientOptions) (*Client, error) {
	a, err := auth.NewClient(opts)
	if err != nil {
		return nil, err
	}

	client := &Client{realmAuth: a}
	err = client.configure(opts)

	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) configure(opts *options.ClientOptions) error {
	if err := opts.Validate(); err != nil {
		return err
	}
	uri := "https://stitch.mongodb.com/api/client/v2.0/app/" + *opts.AppID + "/graphql"
	c.uri = &uri
	return nil
}

//Health needs to be implemented on the GraphQL Server as it looks for a very specific schema/document.
func (c *Client) Health() (interface{}, error) {

	var q struct {
		Health struct {
			ID          string `graphql:"_id"`
			Status      string `graphql:"status"`
			Description string `graphql:"description"`
			Endpoint    string `graphql:"endpoint"`
		} `graphql:"health"`
	}

	err := c.Query(context.TODO(), &q, nil)
	if err != nil {
		log.Fatal(err)
	}

	return q, nil

}

// Connect establishes Realm auth and creates a new graphql client
func (c *Client) Connect() error {
	err := c.realmAuth.Connect()
	if err != nil {
		return err
	}
	c.client = gql.NewClient(*c.uri, c.realmAuth.HttpClient)
	return nil
}

// Disconnect disconnects user session
func (c *Client) Disconnect() error {
	// TODO
	return fmt.Errorf("not yet implemented")
}

// Query wrapper
func (c *Client) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	return c.client.Query(ctx, q, variables)
}

// Mutate wrapper
func (c *Client) Mutate(ctx context.Context, m interface{}, variables map[string]interface{}) error {
	return c.client.Mutate(ctx, m, variables)
}
