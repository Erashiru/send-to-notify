package models

type StarterAppStatus string

const (
	Created      StarterAppStatus = "created"
	Canceled     StarterAppStatus = "canceled"
	Draft        StarterAppStatus = "draft"
	NotConfirmed StarterAppStatus = "notConfirmed"
	Checked      StarterAppStatus = "checked"
	InProgress   StarterAppStatus = "inProgress"
	Cooked       StarterAppStatus = "cooked"
	OnTheWay     StarterAppStatus = "onTheWay"
	Done         StarterAppStatus = "done"
)

func (s StarterAppStatus) String() string {
	return string(s)
}
