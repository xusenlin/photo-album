package dateTime

import "time"

const DateTimeFormat = "2006-01-02 15:04:05"

type DateTime struct {
	time.Time
}
