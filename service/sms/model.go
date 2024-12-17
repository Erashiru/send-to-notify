package sms

type Error struct {
	Code    int
	Message string
}

func (e Error) Error() string {
	return e.Message
}

var (
	CredsErr       = Error{Code: 2, Message: "wrong credential of service"}
	ParamErr       = Error{Code: 1, Message: "invalid parameter"}
	IpErr          = Error{Code: 4, Message: "ip is blocked"}
	RestrictionErr = Error{Code: 5, Message: "sms are blocked to send, due to attempt to send newsletter"}
	MoneyErr       = Error{Code: 3, Message: "not enough money"}
	DateFormatErr  = Error{Code: 6, Message: "invalid date format"}
	InvalidNumErr  = Error{Code: 7, Message: "cannot send sms to this phone number"}
	SpamErr        = Error{Code: 9, Message: "too many requests"}
	SendErr        = Error{Code: 8, Message: "cannot send sms due to sms service error"}
)

type Response struct {
	Error     string `json:"error"`
	ErrorCode int    `json:"error_code"`
	SmsId     int    `json:"id"`
}
