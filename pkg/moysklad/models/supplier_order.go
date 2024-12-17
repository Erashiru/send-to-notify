package models

type SupplierOrder struct {
	Organization MetaData `json:"organization"`
	Agent        MetaData `json:"agent"`
	Description  string   `json:"description,omitempty"`
	//Name                  string       `json:"name,omitempty"`
	//Code                  string       `json:"code,omitempty"`
	//Moment                string       `json:"moment,omitempty"`
	//Applicable            bool         `json:"applicable,omitempty"`
	//VatEnabled            bool         `json:"vatEnabled,omitempty"`
	//VatIncluded           bool         `json:"vatIncluded,omitempty"`
	//State                 MetaData     `json:"state,omitempty"`
	//Store                 MetaData     `json:"store,omitempty"`
	//Contract              MetaData     `json:"contract,omitempty"`
	//Rate                  Rate         `json:"rate,omitempty"`
	//DeliveryPlannedMoment string       `json:"deliveryPlannedMoment,omitempty"`
	//Attributes            []Attributes `json:"attributes,omitempty"`
}
