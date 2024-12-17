package logger

type LoggerInfo struct {
	System   string      `json:"system,omitempty"`
	Request  interface{} `json:"request,omitempty"`
	Response interface{} `json:"response,omitempty"`
	Status   int         `json:"status,omitempty"`
}
