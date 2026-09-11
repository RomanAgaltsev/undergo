// Package drill — C14/05 api-client.
package drill

import (
	"context"
	"encoding/json"
	"net/http"
)

// User is the decoded API response.
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// GetUser fetches a user by URL and decodes it.
func GetUser(ctx context.Context, url string) (User, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return User{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return User{}, err
	}

	var u User
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return User{}, err
	}
	resp.Body.Close()

	return u, nil
}
