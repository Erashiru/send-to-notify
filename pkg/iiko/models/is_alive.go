package models

type IsAliveRequest struct {
	OrganizationIds  []string `json:"organizationIds"`
	TerminalGroupIds []string `json:"terminalGroupIds"`
}

type IsAliveResponse struct {
	CorrelationId string          `json:"correlationId"`
	IsAliveStatus []IsAliveStatus `json:"isAliveStatus"`
}

type IsAliveStatus struct {
	IsAlive         bool   `json:"isAlive"`
	TerminalGroupId string `json:"terminalGroupId"`
	OrganizationId  string `json:"organizationId"`
}
