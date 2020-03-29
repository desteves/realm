// Package auth handles the Realm GraphQL Server authentication.
// This consists of providing valid credentials to obtain a token.
// The token is refreshed every 30 minutes.
package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/desteves/realm/pkg/options"
	"golang.org/x/oauth2"
)

// Client holds a http realm client
type Client struct {
	HTTPClient *http.Client
	Token      *oauth2.Token // public so the application can use withExtra() to access device id or user_id

	//private
	options *options.ClientOptions
	oauth   *oauth2.Config
}

// NewClient creates a new Client with endpoints to Realm based on the provided
// client options.
func NewClient(opts *options.ClientOptions) (*Client, error) {
	client := &Client{options: opts, oauth: &oauth2.Config{}, Token: &oauth2.Token{}}
	err := client.configure(opts)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) configure(opts *options.ClientOptions) error {
	if err := opts.Validate(); err != nil {
		return err
	}
	c.createEndpoint(*opts.AppID, *opts.AuthMechanism)
	return nil
}

func (c *Client) createEndpoint(appid, provider string) {
	c.oauth.Endpoint = oauth2.Endpoint{
		AuthURL:  "https://stitch.mongodb.com/api/client/v2.0/app/" + appid + "/auth/providers/" + provider + "/login",
		TokenURL: "https://stitch.mongodb.com/api/client/v2.0/auth/session",
	}
}

// Ping assumes an http service named "ping" with an incoming_webhook calling a function named "test" which returns 200 has been created.
func (c *Client) Ping() error {
	// TODO - add checks for nil
	uri := "https://webhooks.mongodb-stitch.com/api/client/v2.0/app/" + *c.options.AppID + "/service/ping/incoming_webhook/test"
	resp, err := c.HTTPClient.Get(uri)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response status (%+v)", resp.StatusCode)
	}

	return nil
}

// ConnectWithToken connect to realm with an existing token, either user-provided or internally obtained. Then token needs to be valid for the client to work.
func (c *Client) ConnectWithToken(t *oauth2.Token) error {
	c.HTTPClient = oauth2.NewClient(oauth2.NoContext, c.oauth.TokenSource(oauth2.NoContext, t))
	return nil
}

// Connect connects to realm and establishes http client with auto refresh Token
func (c *Client) Connect() error {
	err := c.retrieveFirstToken()
	if err != nil {
		return err
	}
	return c.ConnectWithToken(c.Token)
}

// func (c *Client) Disconnect() error {
// 	// TODO
// 	// c.HTTPClient.
// 	return fmt.Errorf("not yet implemented")
// }

// Because of non-standard body and headers we need to do a little "hack"
// and request the first Token slightly different than how the
// oauth2 package does it. Need to further explore if we can use the
// native oauth2 functions instead...using this for *now*
func (c *Client) retrieveFirstToken() error {

	b, err := json.Marshal(c.options.Credential)
	if err != nil {
		return err
	}
	body := bytes.NewReader(b)
	req, err := http.NewRequest("POST", c.oauth.Endpoint.AuthURL, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response %+v", resp)
	}
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &c.Token)
	if err != nil {
		return err
	}

	// one day this will be set properly in the response. The docs say 30 mins, so setting it to 29
	// https://docs.mongodb.com/stitch/graphql/authenticate-graphql-requests/#refresh-a-client-api-access-token
	if c.Token.Expiry.IsZero() {
		c.Token.Expiry = time.Now().Add(time.Minute * 29)
	}

	// also storing other "raw" but undocumented fields in the response.
	raw := map[string]interface{}{}
	err = json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}
	delete(raw, "access_token")
	delete(raw, "refresh_token")
	delete(raw, "token_type")
	delete(raw, "expiry")
	c.Token = c.Token.WithExtra(raw)
	return nil
}
