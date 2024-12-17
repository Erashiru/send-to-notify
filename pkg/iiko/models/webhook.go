package models

type GetWebhookSettingRequest struct {
	OrganizationID string `json:"organizationId"`
}

type GetWebhookSettingResponse struct {
	CorrelationID string `json:"correlationId"`
	APILoginName  string `json:"apiLoginName"`
	WebhookURI    string `json:"webHooksUri"`
	AuthToken     string `json:"authToken"`
}

type UpdateWebhookRequest struct {
	OrganizationID string `json:"organizationId"`
	WebhookURI     string `json:"webHooksUri"`
	AuthToken      string `json:"authToken"`
}
