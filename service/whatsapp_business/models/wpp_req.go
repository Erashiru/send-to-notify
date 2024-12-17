package models

type VerificationCodeRequest struct {
	MessagingProduct string       `json:"messaging_product"`
	RecipientType    string       `json:"recipient_type"`
	To               string       `json:"to"`
	Type             string       `json:"type"`
	Template         TemplateData `json:"template"`
}

type TemplateData struct {
	Name       string       `json:"name"`
	Language   LanguageData `json:"language"`
	Components []Component  `json:"components"`
}

type LanguageData struct {
	Code string `json:"code"`
}

type Component struct {
	Type       string      `json:"type"`
	Parameters []Parameter `json:"parameters"`
	SubType    string      `json:"sub_type,omitempty"`
	Index      string      `json:"index,omitempty"`
}

type Parameter struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
