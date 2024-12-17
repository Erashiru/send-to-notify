package http

import (
	"context"
	"fmt"
	httpCli "github.com/kwaaka-team/orders-core/pkg/net-http-client/http"
	"github.com/kwaaka-team/orders-core/pkg/whatsapp/clients"
	"net/url"
)

const baseURL = "https://api.ultramsg.com"

type Client struct {
	Instance  string
	AuthToken string
	BaseUrl   string
	Client    httpCli.Client
	Headers   map[string]string
}

func NewClient(cfg *clients.Config) (clients.Whatsapp, error) {
	if cfg.BaseURL == "" {
		cfg.BaseURL = baseURL
	}
	headers := map[string]string{
		"content-type": "application/x-www-form-urlencoded",
	}

	client := httpCli.NewHTTPClient(cfg.BaseURL)

	cl := Client{
		Instance:  cfg.Instance,
		AuthToken: cfg.AuthToken,
		BaseUrl:   cfg.BaseURL,
		Client:    client,
		Headers:   headers,
	}

	return cl, nil
}

type Req struct {
	Token string `form:"token"`
	To    string `form:"to"`
	Body  string `form:"body"`
}

func (cli Client) SendMessage(ctx context.Context, to, message string) error {
	headers := map[string]string{
		"content-type": "application/x-www-form-urlencoded",
	}

	_, _, err := cli.Client.Post(fmt.Sprintf("/%s/messages/chat?token=%s&to=%s&body=%s&priority=10", cli.Instance, cli.AuthToken, to, message), nil, headers)
	if err != nil {
		return err
	}

	return nil
}

func (cli Client) SendFilePdf(ctx context.Context, to, fileName, message, pdfFileBase64 string) error {
	path := fmt.Sprintf("/%s/messages/document", cli.Instance)
	data := url.Values{}
	data.Set("token", cli.AuthToken)
	data.Set("to", to)
	data.Set("filename", fileName)
	data.Set("document", pdfFileBase64)
	data.Set("caption", message)

	payload := []byte(data.Encode())

	_, _, err := cli.Client.Post(path, payload, cli.Headers)
	if err != nil {
		return err
	}

	return nil
}
