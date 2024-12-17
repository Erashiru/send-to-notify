package models

import "time"

/*
	TODO:
	- when output fare map (fareWithCount map[string]Cost) if there more than 1 fare (of course there would be, because amount of fares
		equals to amount of stores) they would overlap each, not sum, as they should
		- here is additional problem from this: if there are more than 1 fare (probably), besides overlapping problem, there
			would be calculation issue if currency is different (question is - is this my job to calculate the overall amount
			that comes from these fares and if yes, how should I convert one currency to another?)
*/

type LegalEntityForm struct {
	ID               string    `bson:"id,omitempty" json:"id,omitempty"`
	Name             string    `bson:"name" json:"name"`
	BIN              string    `bson:"bin" json:"bin"`
	KNP              string    `bson:"knp" json:"knp"`
	PaymentType      string    `bson:"payment_type" json:"payment_type"`
	LinkedAccManager string    `bson:"linked_acc_manager" json:"linked_acc_manager"`
	SalesID          string    `bson:"sales_id" json:"sales_id"`
	Contacts         []Contact `bson:"contacts" json:"contacts"`
	SalesComment     string    `bson:"sales_comment,omitempty" json:"sales_comment,omitempty"`
	StoreIds         []string  `bson:"store_ids" json:"store_ids"`
	PaymentCycle     int       `bson:"payment_cycle" json:"payment_cycle"`
	Status           string    `bson:"status" json:"status"`
	FirstPaymentAt   time.Time `bson:"first_payment_at,omitempty" json:"first_payment_at,omitempty"`
	CreatedAt        time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt        time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type Contact struct {
	FullName string `bson:"full_name" json:"full_name"`
	Position string `bson:"position" json:"position"`
	Phone    string `bson:"phone" json:"phone"`
	Email    string `bson:"email" json:"email"`
	Comment  string `bson:"comment,omitempty" json:"comment,omitempty"`
}

type LegalEntityView struct {
	ID           string        `bson:"_id" json:"id"`
	Name         string        `bson:"name" json:"name"`
	BIN          string        `bson:"bin" json:"bin"`
	KNP          string        `bson:"knp" json:"knp"`
	PaymentType  string        `bson:"payment_type" json:"payment_type"`
	Contacts     []Contact     `bson:"contacts" json:"contacts"`
	SalesComment string        `bson:"sales_comment,omitempty" json:"sales_comment,omitempty"`
	Status       string        `bson:"status" json:"status"`
	PaymentCycle int           `bson:"payment_cycle" json:"payment_cycle"`
	AccManager   []ContactView `bson:"manager" json:"acc_manager"`
	Sales        []ContactView `bson:"sales" json:"sales"`
	Stores       []Store       `bson:"stores" json:"-"`
	Brands       []Brand       `bson:"brands" json:"-"`
	BusinessInfo BusinessInfo  `bson:"-" json:"business_info"`
}

type Store struct {
	Name         string `bson:"name"`
	WoltExists   bool   `bson:"wolt_exists"`
	GlovoExists  bool   `bson:"glovo_exists"`
	YandexExists bool   `bson:"yandex_exists"`
	Fare         Fare   `bson:"fare"`
}
type Brand struct {
	Name string `bson:"name"`
}

type ContactView struct {
	ID       string `bson:"_id" json:"id"`
	FullName string `bson:"full_name" json:"full_name"`
	Phone    string `bson:"phone" json:"phone"`
}

type BusinessInfo struct {
	Brands            []string     `json:"brands"`
	BrandsCount       int          `json:"brands_count"`
	Stores            []string     `json:"stores"`
	StoresCount       int          `json:"stores_count"`
	Integrations      Integrations `json:"integrations"`
	IntegrationsCount int          `json:"integrations_count"`
	PaymentAmount     Cost         `json:"payment_amount"`
}

type Integrations struct {
	Yandex int `json:"yandex,omitempty"`
	Glovo  int `json:"glovo,omitempty"`
	Wolt   int `json:"wolt,omitempty"`
}

type GetListOfLegalEntitiesDB struct {
	LegalEntityID string    `bson:"_id"`
	Name          string    `bson:"name"`
	PaymentType   string    `bson:"payment_type"`
	Contacts      []Contact `bson:"contacts"`
	Status        string    `bson:"status"`
	Brands        []Brand   `bson:"brands"`
	Stores        []Store   `bson:"stores"`
}

type GetListOfLegalEntities struct {
	LegalEntityID string     `json:"legal_entity_id"`
	Name          string     `json:"name"`
	Brands        []string   `json:"brands"`
	PaymentType   string     `json:"payment_type"`
	Contacts      []Contact  `json:"contacts"`
	PaymentAmount Cost       `json:"payment_amount"`
	Status        string     `json:"status"`
	Fare          FareOutput `json:"fare"`
}

type FareOutput struct {
	Fares       map[string]Fare `json:"fares"`
	OverallCost Cost            `json:"overall_cost"`
}

type Fare struct {
	Type               string `bson:"type" json:"type"`
	Cost               Cost   `bson:"cost" json:"cost"`
	IntegrationsAmount int    `bson:"-" json:"integrations_amount,omitempty"`
}

type Cost struct {
	Value    int    `bson:"value" json:"value"`
	Currency string `bson:"currency" json:"currency"`
}

type Filter struct {
	Search      string
	ContactName string
	PaymentType []string
	Status      []string
}

type GetListOfStoresDB struct {
	Stores []Store `bson:"stores"`
	Brands []Brand `bson:"brands"`
}

type StoreInfo struct {
	Name         string `json:"name"`
	Brand        string `json:"brand"`
	Fare         Fare   `json:"fare"`
	WoltExists   bool   `json:"wolt_exists"`
	GlovoExists  bool   `json:"glovo_exists"`
	YandexExists bool   `json:"yandex_exists"`
}

type GetListOfStores struct {
	StoresInfo     []StoreInfo `json:"store_info"`
	OverallPayment Cost        `json:"overall_payment"`
}
