package license

import (
	"log"
	"strings"

	"github.com/vigneshuvi/GoDateFormat"
	"github.com/afeldman/go-util/time"
)

func GetLicense(lice, email, author, project string) (string, error) {

	lice, err := FetchLicense(lice)
	if err != nil {
		log.Println(err)
		return "", err
	}

	today := time_util.GetToday(GoDateFormat.ConvertFormat("YYYY"))
	lice = strings.Replace(lice, "[year]", today, 1)
	lice = strings.Replace(lice, "[fullname]", author, 1)
	lice = strings.Replace(lice, "[email]", email, 1)
	lice = strings.Replace(lice, "[project]", project, 1)

	return lice, nil
}
