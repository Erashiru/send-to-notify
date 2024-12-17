package models

import "time"

type EntityChangesHistory struct {
	Author           string    `bson:"author"`
	OperationType    string    `bson:"operation_type"`
	CallFunction     string    `bson:"call_function"`
	TaskType         string    `bson:"task_type"`
	RepositoryMethod string    `bson:"repository_method"`
	CollectionName   string    `bson:"collection_name"`
	OldBody          any       `bson:"body"`
	ModifiedAt       time.Time `bson:"modified_at"`
}

type EntityChangesHistoryRequest struct {
	Author   string `json:"author"`
	TaskType string `json:"task_type"`
}
