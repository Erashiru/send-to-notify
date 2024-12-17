package models

import (
	"encoding/json"
	"strings"
	"time"
)

type DateTime time.Time

func (dt DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(dt.String())
}

func (dt *DateTime) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	str = strings.Replace(str, " ", "T", 1)

	str = dt.AddZone(str)

	if strings.Contains(str, ".") {
		runes := []rune(str)
		runes = append(runes[:strings.Index(str, ".")], 'Z')
		str = string(runes)
	}
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}
	*dt = DateTime(t)
	return nil
}

func (dt DateTime) String() string {
	return time.Time(dt).Format(time.RFC3339)
}

func (dt DateTime) AddZone(str string) string {
	if len(str) > 0 && str[len(str)-1] != 'Z' {
		str += "Z"
	}
	return str
}
