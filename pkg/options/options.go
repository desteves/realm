// Package options contains shared properties/helpers to be consumed by other packages.
package options

import "fmt"

// ClientOptions to connect to Realm
type ClientOptions struct {
	AppID         *string     `yaml:"appid" json:"app_id,omitempty"`
	AuthMechanism *string     `yaml:"provider" json:"provider,omitempty"`
	Credential    *Credential `yaml:"credential,omitempty" json:"credential,omitempty"`
}

// Credential are provider-agnostic, fill only needed or omit if using anonymous authentication
type Credential struct {
	Username *string `json:"username,omitempty" yaml:"username,omitempty"`
	Password *string `json:"password,omitempty" yaml:"password,omitempty"`
	Key      *string `json:"key,omitempty" yaml:"key,omitempty"`
	Token    *string `json:"token,omitempty" yaml:"token,omitempty"`
}

// Validate validates the client options. This method will return the first error found.
func (c *ClientOptions) Validate() error {

	if c.AppID == nil {
		return fmt.Errorf("AppID is required, but missing")
	}
	if c.AuthMechanism == nil {
		return fmt.Errorf("Auth Provider is required, but missing")
	}

	switch *(c.AuthMechanism) {
	case "anon-user":
		c.Credential = nil
	case "local-userpass":
		if c.Credential.Username == nil || c.Credential.Password == nil || c.Credential.Key != nil || c.Credential.Token != nil {
			return fmt.Errorf("Wrong creds format for local-userpass")
		}
	case "oauth2-google":
		if c.Credential.Username != nil || c.Credential.Password != nil || c.Credential.Key == nil || c.Credential.Token != nil {
			return fmt.Errorf("Wrong creds format for oauth2-google")
		}

	case "key":
		if c.Credential.Username != nil || c.Credential.Password != nil || c.Credential.Key == nil || c.Credential.Token != nil {
			return fmt.Errorf("Wrong creds format for key")
		}
	case "custom-token":
		if c.Credential.Username != nil || c.Credential.Password != nil || c.Credential.Key != nil || c.Credential.Token == nil {
			return fmt.Errorf("Wrong creds format for custom-token")
		}
	// case "jwtTokenString" :

	default:
		return fmt.Errorf("Auth Provider is not supported")
	}
	return nil
}
