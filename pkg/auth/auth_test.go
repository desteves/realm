package auth

import (
	"testing"

	"github.com/desteves/realm/pkg/options"
)

// Test public funcs
func TestNewClient(t *testing.T) {
	t.Log("Creating a new Realm Client with anonymous authentication ... ( expected err: nil )")
	appid := "graphqlserver-lrnqt"
	// TODO add other mechanisms
	auths := [...]string{"anon-user"}
	for _, auth := range auths {
		opts := &options.ClientOptions{
			AppID:         &appid,
			AuthMechanism: &auth}
		_, err := NewClient(opts)
		if err != nil {
			t.Error("Failed to obtain a new client with " + auth + " auth.")
		}
		t.Log("Got a realm client with " + auth + " auth.")
	}
}

func TestPing(t *testing.T) {

}

func TestConnect(t *testing.T) {

}

// Test private funcs

func TestRetrieveFirstToken(t *testing.T) {

	// unmarshall error, invalid json
	// invalid creds
	// valid creds -- for each auth provider
}
