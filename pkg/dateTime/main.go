package dateTime

import (
	"time"
)

const DateFormat = "2006-01-02"
const DateTimeFormat = "2006-01-02 15:04:05"

type DateTime struct {
	time.Time
}

func (ct *DateTime) MarshalYAML() (interface{}, error) {
	return ct.Time.Format(DateFormat), nil
}

func (ct *DateTime) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var formattedTime string
	if err := unmarshal(&formattedTime); err != nil {
		return err
	}

	parsedTime, err := time.Parse(DateFormat, formattedTime)
	if err != nil {
		return err
	}

	ct.Time = parsedTime
	return nil
}
