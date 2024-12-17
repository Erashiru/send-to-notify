package dto

import (
	"encoding/json"
	"time"
)

const (
	MessageReceived string = "message_received"
)

type WebhookEvent struct {
	EventType   string      `json:"event_type"`
	InstanceId  string      `json:"instance_id"`
	Id          string      `json:"id"`
	ReferenceId string      `json:"reference_id"`
	Data        WebhookData `json:"data"`
}

type WebhookData struct {
	Id           string   `json:"id"`
	From         string   `json:"from"`
	To           string   `json:"to"`
	Author       string   `json:"author"`
	Pushname     string   `json:"pushname"`
	Ack          string   `json:"ack"`
	Type         string   `json:"type"`
	Body         string   `json:"body"`
	Media        string   `json:"media"`
	FromMe       bool     `json:"fromMe"`
	Self         bool     `json:"self"`
	IsForwarded  bool     `json:"isForwarded"`
	IsMentioned  bool     `json:"isMentioned"`
	QuoteMsg     QuoteMsg `json:"quotedMsg"`
	MentionedIds []string `json:"mentionedIds"`
	Time         UnixTime `json:"time"`
}

type UnixTime struct {
	time.Time
}

func (t *UnixTime) UnmarshalJSON(b []byte) error {
	var unixTime int64
	if err := json.Unmarshal(b, &unixTime); err != nil {
		return err
	}
	t.Time = time.Unix(unixTime, 0)
	return nil
}

func (t UnixTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Unix())
}

type QuoteMsg struct {
	Id   string `json:"id"`
	Body string `json:"body"`
}
