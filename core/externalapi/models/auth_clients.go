package models

import (
	"github.com/kwaaka-team/orders-core/pkg/iiko/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	YANDEX = "yandex"
	ADMIN  = "admin"
	EMENU  = "emenu"
)

type SetCredsToStore struct {
	RestID        string              `json:"restaurant_id"`
	AuthId        string              `json:"auth_id"`
	StoreID       []string            `json:"store_id"`
	MenuUrl       string              `json:"menu_url"`
	SendToPos     bool                `json:"send_to_pos"`
	IsMarketplace bool                `json:"is_marketplace"`
	PaymentTypes  models.PaymentTypes `json:"payment_types"`
}

type AuthClient struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Restaurants      []string           `bson:"restaurants"`
	ExternalStoreIDs []string           `bson:"external_store_ids"`
	ClientID         string             `bson:"client_id"`
	ClientSecret     string             `bson:"client_secret"`
	GrantType        string             `bson:"grant_type"`
	Scope            string             `bson:"scope"`
	Service          string             `bson:"service"`
	ExpirationDate   time.Time          `bson:"expiration_date"`
	GrantedBy        string             `bson:"granted_by"`
	LastUsedDate     time.Time          `bson:"last_used_date"`
	CreatedAt        time.Time          `bson:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at"`
}

type AuthenticateData struct {
	ClientID     string `json:"client_id" validate:"required" form:"client_id"`
	ClientSecret string `json:"client_secret" validate:"required" form:"client_secret"`
	GrantType    string `json:"grant_type" validate:"required" form:"grant_type"`
	Scope        string `json:"scope" validate:"required" form:"scope"`
}

type Credentials struct {
	RestID           string `json:"restaurant_id"`
	AuthenticateData `json:"authenticate_data"`
	Service          string `json:"service"`
}

func ToModelAuthClient(y Credentials) AuthClient {
	return AuthClient{
		ClientID:     y.AuthenticateData.ClientID,
		ClientSecret: y.AuthenticateData.ClientSecret,
		Scope:        y.AuthenticateData.Scope,
		GrantType:    y.AuthenticateData.GrantType,
		UpdatedAt:    time.Now().UTC(),
		CreatedAt:    time.Now().UTC(),
		GrantedBy:    ADMIN,
		Service:      y.Service,
	}
}
