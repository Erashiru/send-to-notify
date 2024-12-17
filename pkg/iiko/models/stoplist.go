package models

type StopListRequest struct {
	Organizations []string `json:"organizationIds"`
}

type StopListResponse struct {
	TerminalGroups []TerminalGroup `json:"terminalGroupStopLists"`
}

type TerminalGroup struct {
	Organization string         `json:"organizationId"`
	Items        []TerminalItem `json:"items"`
}

func (s StopListResponse) Item(terminalID string) (TerminalItem, error) {

	for i := range s.TerminalGroups {
		for _, terminalItem := range s.TerminalGroups[i].Items {
			if terminalItem.TerminalGroupID == terminalID {
				return terminalItem, nil
			}
		}
	}

	return TerminalItem{}, nil
}

type TerminalInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Addr string `json:"address"`
}

type TerminalItem struct {
	TerminalGroupID string         `json:"terminalGroupId"`
	Items           []StopListItem `json:"items,omitempty"`
}

type StopListItem struct {
	Balance   float64 `json:"balance"`
	ProductID string  `json:"productId"`
}
