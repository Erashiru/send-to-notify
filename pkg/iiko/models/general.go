package models

import (
	"encoding/json"
	"reflect"
	"strings"
)

type Type int

const (
	DISH Type = iota + 1
	GOOD
	MODIFIER
	SERVICE
)

var types = map[Type]string{
	DISH:     "dish",
	GOOD:     "good",
	MODIFIER: "modifier",
	SERVICE:  "service",
}

func (t Type) MarshalJSON() ([]byte, error) {
	str := types[t]
	return json.Marshal(str)
}

func (t *Type) UnmarshalJSON(data []byte) error {
	var str string
	if string(data) == `null` {
		return nil
	}
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	for k, v := range types {
		if v == strings.ToLower(str) {
			*t = k
			return nil
		}
	}

	return &json.InvalidUnmarshalError{
		Type: reflect.TypeOf(t),
	}
}

type CorID struct {
	CorID string `json:"correlationId"`
}

type OrganizationRequest struct {
	OrganizationID string `json:"organizationId"`
}

type AdditionalInfo struct {
	ReturnAddInfo bool `json:"returnAdditionalInfo"`
}
