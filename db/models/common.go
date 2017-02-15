package models

import (
	"time"
	"strings"
)

type JSONTime struct {
	time.Time
}

func (t JSONTime)MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Time.Format(time.RFC1123) + `"`), nil
}

func (t *JSONTime)UnmarshalJSON(b []byte) error {
	tt, err := time.Parse(time.RFC1123, strings.Trim(string(b), `"`))
    	if err != nil {
		return err
    	}
    	t.Time = tt
    	return nil
}