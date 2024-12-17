package poster

import (
	"encoding/json"
	"fmt"
)

type WHEventDto struct {
	Account       string      `json:"account" example:"test"`
	AccountNumber string      `json:"account_number" example:"526209"`
	Object        string      `json:"object" example:"transaction"`
	ObjectID      json.Number `json:"object_id" example:"24876"`
	Action        string      `json:"action" example:"changed"`
	Time          string      `json:"time" example:"1695621950"`
	Verify        string      `json:"verify" example:"faaac2fe811509fd82fe83be47ab5c52"`
	Data          string      `json:"data" example:"{\"type\":1}"`
}

func (u *WHEvent) UnmarshalJSON(data []byte) error {
	var tmp WHEventDto
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	u.fromDto(tmp)

	return nil
}

func (u *WHEvent) fromDto(r WHEventDto) {
	u.ObjectID = string(r.ObjectID)
	u.AccountNumber = r.AccountNumber
	u.Account = r.Account
	u.Object = r.Object
	u.Action = r.Action
	u.Time = r.Time
	u.Verify = r.Verify
	u.Data = r.Data
}

type WHEvent struct {
	Account       string `json:"account" example:"test"`
	AccountNumber string `json:"account_number" example:"526209"`
	Object        string `json:"object" example:"transaction"`
	ObjectID      string `json:"object_id" example:"24876"`
	Action        string `json:"action" example:"changed"`
	Time          string `json:"time" example:"1695621950"`
	Verify        string `json:"verify" example:"faaac2fe811509fd82fe83be47ab5c52"`
	Data          string `json:"data" example:"{\"type\":1}"`
}

type Data struct {
	Type          int     `json:"type,omitempty"`
	ElementId     int     `json:"element_id,omitempty"`
	StorageId     int     `json:"storage_id,omitempty"`
	ValueRelative float64 `json:"value_relative,omitempty"`
	ValueAbsolute float64 `json:"value_absolute,omitempty"`
	ProductId     string  `json:"product_id,omitempty"`
}

func CastData(data string) (Data, error) {
	var response Data
	err := json.Unmarshal([]byte(data), &response)
	if err != nil {
		fmt.Println("Error:", err)
		return Data{}, err
	}
	return response, nil
}

type RestaurantStoplistItems struct {
	ProductID   string
	AttributeID string
	IsAvailable bool
}
