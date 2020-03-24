package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/desteves/realm/internal/gqlquery"
	"github.com/desteves/realm/pkg/auth"
	"github.com/desteves/realm/pkg/options"
	"github.com/pkg/errors"
	"golang.org/x/net/context/ctxhttp"
)

//HealthCheck Query for Realm Config Verification
type HealthCheck struct {
	ID          string
	Description string
	Status      string
	Endpoint    string
}

// PathSegment returned when the response had an error. Path segments that represent fields should be strings, and path segments that represent list indices should be 0‐indexed integers. If the error happens in an aliased field, the path to the error should use the aliased name, since it represents a path in the response, not in the query.
type PathSegment struct {
	parent *PathSegment
	value  interface{}
}

func (p *PathSegment) toSlice() []interface{} {
	if p == nil {
		return nil
	}
	return append(p.parent.toSlice(), p.value)
}

// Location returned when the response had an error. Each location is a map with the keys line and column, both positive numbers starting from 1 which describe the beginning of an associated syntax element
type Location struct {
	Line   int `json:"line,omitempty"`
	Column int `json:"column,omitempty"`
}

// Error is a GraphQL Error as per http://spec.graphql.org/June2018/#sec-Errors
type Error struct {
	Message    string                 `json:"message,omitempty"`
	Path       PathSegment            `json:"path,omitempty"`
	Locations  []Location             `json:"locations,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type Variable map[string]interface{}

// Request is the payload for GraphQL
type Request struct {
	Query         string   `json:"query" graphql:"query"`
	Variables     Variable `json:"variables" graphql:"variables"`
	OperationName *string  `json:"operationName" graphql:"operationName"`
}

// Response is the payload for a GraphQL response.
type Response struct {
	Data   interface{} `json:"data,omitempty" graphql:"data,omitempty"`
	Errors []Error     `json:"errors,omitempty"`
}

// Client is a Realm GraphQL Client with authentication to a Realm Application
type Client struct {
	client *auth.Client
	uri    *string
}

// NewClient creates a new Client
func NewClient(opts *options.ClientOptions) (*Client, error) {
	a, err := auth.NewClient(opts)
	if err != nil {
		return nil, err
	}

	client := &Client{client: a}
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
func (c *Client) Health(response *Response) error {
	var q struct {
		Health struct {
			ID          string `graphql:"_id"`
			Status      string `graphql:"status"`
			Description string `graphql:"description"`
			Endpoint    string `graphql:"endpoint"`
		} `graphql:"health"`
	}
	return c.Query(context.TODO(), &q, nil, response)
}

// Connect establishes Realm auth and creates a new graphql client
func (c *Client) Connect() error {
	err := c.client.Connect()
	if err != nil {
		return err
	}

	return nil
}

// // Disconnect disconnects user session
// func (c *Client) Disconnect() error {
// 	return c.client.Disconnect()
// }

// Query wrapper
func (c *Client) Query(ctx context.Context, query interface{}, variables map[string]interface{}, response *Response) error {
	payload := Request{
		Query:     gqlquery.ConstructQuery(query, variables),
		Variables: variables,
	}

	err := c.do(ctx, payload, response)
	if err != nil {
		return errors.Wrap(err, "q do")
	}

	return nil
}

// Mutate wrapper
func (c *Client) Mutate(ctx context.Context, mutation string, variables map[string]interface{}, response *Response) error {
	payload := Request{
		Query:     mutation,
		Variables: variables,
	}
	// fmt.Printf("\n Payload is :\n %+v \n", payload)
	err := c.do(ctx, payload, response)
	if err != nil {
		return errors.Wrap(err, "m do")
	}

	return nil
}

// do executes a single GraphQL operation.
func (c *Client) do(ctx context.Context, payload interface{}, response *Response) error {

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		return err
	}
	resp, err := ctxhttp.Post(ctx, c.client.HttpClient, *c.uri, "application/json", &buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("non-200 OK status code: %v body: %q", resp.Status, body)
	}

	return json.NewDecoder(resp.Body).Decode(&response)
}
