package license

import (
	"log"
	"strings"
	"time"

	"github.com/vigneshuvi/GoDateFormat"
)

func GetLicense(lice, email, author, project string) string {

	lice, err := FetchLicense(lice)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	//now := time.Now().Year()
	today := GetToday(GoDateFormat.ConvertFormat("YYYY"))
	lice = strings.Replace(lice, "[year]", today, 1)
	lice = strings.Replace(lice, "[fullname]", author, 1)
	lice = strings.Replace(lice, "[email]", email, 1)
	lice = strings.Replace(lice, "[project]", project, 1)

	return lice
}

func GetToday(format string) (todayString string) {
	today := time.Now()
	todayString = today.Format(format)
	return
}
