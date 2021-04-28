package controllersv1

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

const (
	tokenBearer = "JEDI-EHYJzdWIiOiJkZmRmc2RmZHMiLCJuYW1lIjP0"
	apiUrl      = "http://localhost:8001"
)

func TestGetSpacecrafts(t *testing.T) {
	endpoint := "/api/v1/spacecrafts"

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = endpoint
	urlStr := u.String()

	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodGet, urlStr, nil) // URL-encoded payload
	r.Header.Add("Authorization", "Bearer "+tokenBearer)

	resp, _ := client.Do(r)
	if resp.StatusCode != http.StatusOK {
		t.Error(fmt.Sprintf("The Status must be 200 not %s", resp.Status))
	}
}
