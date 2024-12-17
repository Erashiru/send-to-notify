package models

type TerminalGroupsResponse struct {
	TerminalGroups []TerminalGroupInfo `json:"terminalGroups"`
}

type TerminalGroupInfo struct {
	Organization string             `json:"organizationId"`
	Items        []TerminalItemInfo `json:"items"`
}

type TerminalItemInfo struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organizationId"`
}
