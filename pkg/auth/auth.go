package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/desteves/realm/pkg/options"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// Client holds a http real client
type Client struct {
	HttpClient *http.Client
	options    *options.ClientOptions
	oauth      *oauth2.Config
	Token      *oauth2.Token // the application can use withExtra() to access device id or user_id
}

// NewClient creates a new Client
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

	uri := "https://webhooks.mongodb-stitch.com/api/client/v2.0/app/" + *c.options.AppID + "/service/ping/incoming_webhook/test"
	resp, err := c.HttpClient.Get(uri)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response status (%+v)", resp.StatusCode)
	}

	return nil
}

// Connect connects to realm and establishes http client with auto refresh Token
func (c *Client) Connect() error {
	err := c.retrieveFirstToken()
	if err != nil {
		return err
	}
	c.HttpClient = oauth2.NewClient(oauth2.NoContext, c.oauth.TokenSource(oauth2.NoContext, c.Token))
	return nil
}

func (c *Client) Disconnect() error {

	// TODO
	return fmt.Errorf("not yet implemented")
}

// Because of non-standard body and headers we need to do a little "hack"
// and request the first Token slightly different than how the
// oauth2 package does it. Need to further explore if we can use the
// native oauth2 functions instead...using this for *now*
func (c *Client) retrieveFirstToken() error {

	b, err := json.Marshal(c.options.Credential)
	if err != nil {
		log.Fatal(err)
	}
	body := bytes.NewReader(b)
	req, err := http.NewRequest("POST", c.oauth.Endpoint.AuthURL, body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("! bad response %+v", resp)
	}
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(b, &c.Token)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
