package license

import (
	"strings"
)

func findPlaceholders(body string, keys []string) (folders []string) {

	for _, k := range keys {
		if strings.Contains(body, k) {
			folders = append(folders, k)
		}
	}
	return
}
