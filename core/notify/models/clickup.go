package models

type ClickUpTask struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Assignees   []int  `json:"assignees"`
	Status      string `json:"status"`
	Priority    int    `json:"priority"`
	DueDate     int64  `json:"due_date"`
	DueDateTime bool   `json:"due_date_time"`
}

type ClickUpTaskResponse struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}
