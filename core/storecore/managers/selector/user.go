package selector

type User struct {
	ID               string   `bson:"_id,omitempty"`
	Username         string   `bson:"username"`
	StoreId          string   `bson:"restaurant_id"`
	StoreGroupId     string   `bson:"restaurant_group_id"`
	FcmTokens        []string `bson:"fcm_tokens"`
	SendNotification bool     `bson:"send_notification"`
}

func NewEmptyUserSearch() User {
	return User{}
}

func (u *User) HasUsername() bool {
	return u.Username != ""
}

func (u *User) HasID() bool {
	return u.ID != ""
}

func (u *User) HasStoreID() bool {
	return u.StoreId != ""
}

func (u *User) HasStoreGroupId() bool {
	return u.StoreGroupId != ""
}

func (u User) SetID(id string) User {
	u.ID = id
	return u
}

func (u User) SetSendNotification(sendNotification bool) User {
	u.SendNotification = sendNotification
	return u
}

func (u User) SetUsername(username string) User {
	u.Username = username
	return u
}

func (u User) SetStoreID(id string) User {
	u.StoreId = id
	return u
}

func (u User) SetStoreGroupID(id string) User {
	u.StoreGroupId = id
	return u
}
