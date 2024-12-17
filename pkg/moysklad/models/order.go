package models

type Order struct {
	Context Context `json:"context"`
	Meta    Meta    `json:"meta"`
	Rows    []Rows  `json:"rows"`
}

type Context struct {
	Employee MetaData `json:"employee"`
}

type Meta struct {
	Href         string  `json:"href"`
	MetadataHref *string `json:"metadataHref,omitempty"`
	Type         string  `json:"type"`
	MediaType    string  `json:"mediaType,omitempty"`
	Size         *int    `json:"size,omitempty"`
	Limit        *int    `json:"limit,omitempty"`
	Offset       *int    `json:"offset,omitempty"`
}

type Rows struct {
	Meta                Meta         `json:"meta"`
	ID                  string       `json:"id"`
	AccountID           string       `json:"accountId"`
	SyncID              string       `json:"syncId"`
	Name                string       `json:"name"`
	Description         string       `json:"description"`
	ExternalCode        string       `json:"externalCode"`
	Owner               MetaData     `json:"owner"`
	Shared              bool         `json:"shared"`
	Group               MetaData     `json:"group"`
	Printed             bool         `json:"printed"`
	Published           bool         `json:"published"`
	VatEnabled          bool         `json:"vatEnabled"`
	VatIncluded         bool         `json:"vatIncluded"`
	Sum                 float64      `json:"sum"`
	ReservedSum         float64      `json:"reservedSum"`
	PayedSum            float64      `json:"payedSum"`
	ShippedSum          float64      `json:"shippedSum"`
	InvoicedSum         float64      `json:"invoicedSum"`
	TaxSystem           string       `json:"taxSystem"`
	Rate                Rate         `json:"rate"`
	Organization        MetaData     `json:"organization"`
	Store               MetaData     `json:"store"`
	Contract            MetaData     `json:"contract"`
	Agent               MetaData     `json:"agent"`
	State               MetaData     `json:"state"`
	OrganizationAccount MetaData     `json:"organizationAccount"`
	AgentAccount        MetaData     `json:"agentAccount"`
	SalesChannel        MetaData     `json:"salesChannel"`
	Attributes          []Attributes `json:"attributes"`
	Positions           MetaData     `json:"positions"`

	ShipmentAddress string  `json:"shipmentAddress"`
	Address         Address `json:"shipmentAddressFull"`

	DeliveredPlannedMoment string `json:"deliveredPlannedMoment"`
	Moment                 string `json:"moment"`
	Created                string `json:"created"`
	Updated                string `json:"updated"`
}

type MetaData struct {
	Meta Meta `json:"meta"`
}

type Rate struct {
	Currency MetaData `json:"currency"`
	Value    int      `json:"value"`
}

type Attributes struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Value float64 `json:"value"`
	Meta  Meta    `json:"meta"`
}

type Address struct {
	PostalCode string   `json:"postalCode"`
	Country    MetaData `json:"country"`
	Region     MetaData `json:"region"`
	City       string   `json:"city"`
	Street     string   `json:"street"`
	House      string   `json:"house"`
	Apartment  string   `json:"apartment"`
	AddInfo    string   `json:"addInfo"`
	Comment    string   `json:"comment"`
}
