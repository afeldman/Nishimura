package license

import (
	"context"
	"fmt"
	"net/http"

	"strings"

	"github.com/google/go-github/github"
)

func fetchLicense(key string) (string, error) {

	keq = string.ToLower(key)

	// Create default client
	client := github.NewClient(nil)

	// Fetch a LICENSE from Github API
	Debugf("Fetch license from GitHub API by key: %s", key)
	license, res, err := client.Licenses.Get(context.Background(), key)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid status code from GitHub\n %s\n", res.String())
	}
	Debugf("Fetched license name: %s", *license.Name)

	return *license.Body, nil
}
