package dto

type BitrixEventDataRequest struct {
	Event        string              `form:"event" json:"event"`
	EventID      string              `form:"event_handler_id" json:"event_handler_id"`
	Timestamp    string              `form:"ts" json:"timestamp"`
	DataFieldsID string              `form:"data[FIELDS][ID]" json:"data_fields_id"`
	Auth         BitrixEventDataAuth `form:"auth" json:"auth"`
}

type BitrixEventDataAuth struct {
	Domain           string `form:"auth[domain]" json:"domain"`
	ClientEndpoint   string `form:"auth[client_endpoint]" json:"client_endpoint"`
	ServerEndpoint   string `form:"auth[server_endpoint]" json:"server_endpoint"`
	MemberID         string `form:"auth[member_id]" json:"member_id"`
	ApplicationToken string `form:"auth[application_token]" json:"application_token"`
}
