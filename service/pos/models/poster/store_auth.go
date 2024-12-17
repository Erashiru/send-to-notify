package poster

import "time"

type PosterStoreAuthError struct {
	Code         int    `json:"code"`
	ErrorType    string `json:"error_type"`
	ErrorMessage string `json:"error_message"`
}

type PosterStoreAuth struct {
	ID            string    `bson:"_id,omitempty"`
	AccessToken   string    `json:"access_token" bson:"access_token"`
	AccountNumber string    `json:"account_number" bson:"account_number"`
	User          User      `json:"user" bson:"user"`
	OwnerInfo     OwnerInfo `json:"ownerInfo" bson:"owner_info"`
	CreatedAt     time.Time `bson:"created_at"`
}
type User struct {
	ID     int    `json:"id" bson:"id"`
	Name   string `json:"name" bson:"name"`
	Email  string `json:"email" bson:"email"`
	RoleID int    `json:"role_id" bson:"role_id"`
}
type OwnerInfo struct {
	Email       string `json:"email" bson:"email"`
	Phone       string `json:"phone" bson:"phone"`
	City        string `json:"city" bson:"city"`
	Country     string `json:"country" bson:"country"`
	Name        string `json:"name" bson:"name"`
	CompanyName string `json:"company_name" bson:"company_name"`
}
