package clients

import (
	"context"
)

type Config struct {
	Protocol  string
	BaseURL   string
	Insecure  bool
	Instance  string
	AuthToken string
}

type Whatsapp interface {
	SendMessage(ctx context.Context, to, message string) error
	SendFilePdf(ctx context.Context, to, fileName, message, pdfFileBase64 string) error
}
