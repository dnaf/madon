package gondole

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sendgrid/rest"
)

// UserToken represents a user token as returned by the Mastodon API
type UserToken struct {
	AccessToken string `json:"access_token"`
	CreatedAt   int    `json:"created_at"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

// LoginBasic does basic user authentication
func (g *Client) LoginBasic(username, password string, scopes []string) error {
	if username == "" {
		return fmt.Errorf("missing username")
	}
	if password == "" {
		return fmt.Errorf("missing password")
	}

	hdrs := make(map[string]string)
	opts := make(map[string]string)

	hdrs["User-Agent"] = "Gondole/" + GondoleVersion

	opts["grant_type"] = "password"
	opts["client_id"] = g.ID
	opts["client_secret"] = g.Secret
	opts["username"] = username
	opts["password"] = password
	if len(scopes) > 0 {
		opts["scope"] = strings.Join(scopes, " ")
	}

	req := rest.Request{
		BaseURL:     g.InstanceURL + "/oauth/token",
		Headers:     hdrs,
		QueryParams: opts,
		Method:      rest.Post,
	}

	r, err := rest.API(req)
	if err != nil {
		return err
	}

	var resp UserToken

	err = json.Unmarshal([]byte(r.Body), &resp)
	if err != nil {
		return err
	}

	g.userToken = &resp
	return nil
}
