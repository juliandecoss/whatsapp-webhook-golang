package utils

import "time"

func GetCurrentTime() string {
	location, err := time.LoadLocation("America/Mexico_City")
	if err != nil {
		panic(err)
	}

	now := time.Now().In(location)

	return now.Format("02/01/2006 15:04")
}
