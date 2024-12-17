package selector

type User struct {
	UID         string
	Name        string
	PhoneNumber string
}

func NewEmptyUser() User {
	return User{}
}

func (u User) SetUID(uid string) User {
	u.UID = uid
	return u
}

func (u User) HasUID() bool {
	return u.UID != ""
}

func (u User) SetName(name string) User {
	u.Name = name
	return u
}

func (u User) HasName() bool {
	return u.Name != ""
}

func (u User) SetPhoneNumber(number string) User {
	u.PhoneNumber = number
	return u
}

func (u User) HasPhoneNumber() bool {
	return u.PhoneNumber != ""
}
