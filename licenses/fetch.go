package license

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
)

func FetchLicense(key string) (string, error) {

	key = strings.ToLower(key)

	// Create default client
	client := github.NewClient(nil)

	// Fetch a LICENSE from Github API
	log.Printf("Fetch license from GitHub API by key: %s\n", key)
	license, res, err := client.Licenses.Get(context.Background(), key)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid status code from GitHub\n %s\n", res.String())
	}
	log.Printf("Fetched license name: %s\n", *license.Name)

	return *license.Body, nil
}
