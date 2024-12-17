package whatsapp

type Message struct {
	CustomerPhone string `json:"customer_phone"`
	Message       string `json:"message"`
	InstanceId    string `json:"instance_id"`
	AuthToken     string `json:"auth_token"`
}
