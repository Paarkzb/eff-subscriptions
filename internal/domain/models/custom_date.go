package models

import (
	"encoding/json"
	"strings"
	"time"
)

type CustomDate time.Time

func (cd *CustomDate) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), "\"")
	t, err := time.Parse("01-2006", s)
	if err != nil {
		return err
	}
	*cd = CustomDate(t)
	return nil
}

func (cd *CustomDate) MarshalJSON() ([]byte, error) {
	t := time.Time(*cd)
	return json.Marshal(t.Format("01-2006"))
}

func (cd *CustomDate) Time() time.Time {
	return time.Time(*cd)
}
