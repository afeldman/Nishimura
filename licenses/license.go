package license

import (
	"log"
	"strings"
	"time"

	"github.com/vigneshuvi/GoDateFormat"
)

func GetLicense(lice, email, author, project string) (string, error) {

	lice, err := FetchLicense(lice)
	if err != nil {
		log.Println(err)
		return "", err
	}

	today := GetToday(GoDateFormat.ConvertFormat("YYYY"))
	lice = strings.Replace(lice, "[year]", today, 1)
	lice = strings.Replace(lice, "[fullname]", author, 1)
	lice = strings.Replace(lice, "[email]", email, 1)
	lice = strings.Replace(lice, "[project]", project, 1)

	return lice, nil
}

func GetToday(format string) (todayString string) {
	today := time.Now()
	todayString = today.Format(format)
	return
}
