package http

import (
	"bytes"
	"context"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"net/url"
	"time"
)

type cli struct {
	client  *http.Client
	baseURL *url.URL
}

func NewHTTPClient(baseURL string) Client {
	return &cli{
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(&url.URL{
					Scheme: "http",
					User:   url.UserPassword("max0Q8ga", "mAxsWX6F2Y"),
					Host:   "91.147.127.252:50100",
				}),
			},
		},
		baseURL: &url.URL{
			Host: baseURL,
		},
	}
}

type Client interface {
	Get(endpoint string, body []byte, headers Headers) (int, []byte, error)
	Post(endpoint string, body []byte, headers Headers) (int, []byte, error)
	Patch(endpoint string, body []byte, headers Headers) (int, []byte, error)
	Put(endpoint string, body []byte, headers Headers) (int, []byte, error)
	Delete(endpoint string, body []byte, headers Headers) (int, []byte, error)
}

const defaultHTTPTimeout = time.Minute * 1

type Headers map[string]string

func (cli cli) Post(endpoint string, body []byte, headers Headers) (int, []byte, error) {
	return cli.request(http.MethodPost, endpoint, body, headers)
}

func (cli cli) Get(endpoint string, body []byte, headers Headers) (int, []byte, error) {
	return cli.request(http.MethodGet, endpoint, body, headers)
}

func (cli cli) Patch(endpoint string, body []byte, headers Headers) (int, []byte, error) {
	return cli.request(http.MethodPatch, endpoint, body, headers)
}

func (cli cli) Put(endpoint string, body []byte, headers Headers) (int, []byte, error) {
	return cli.request(http.MethodPut, endpoint, body, headers)
}

func (cli cli) Delete(endpoint string, body []byte, headers Headers) (int, []byte, error) {
	return cli.request(http.MethodDelete, endpoint, body, headers)
}

func (cli cli) request(method, endpoint string, body []byte, headers Headers) (int, []byte, error) {
	path := cli.baseURL.Host + endpoint

	log.Info().Msgf("request path: %s", path)
	log.Info().Msgf("request body: %s", string(body))

	req, err := http.NewRequest(method, path, bytes.NewReader(body))
	if err != nil {
		return 0, nil, err
	}

	for header, value := range headers {
		req.Header.Set(header, value)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultHTTPTimeout)
	defer cancel()

	response, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return 0, nil, err
	}

	return extractCodeAndBodyFromResponse(response)
}

func extractCodeAndBodyFromResponse(response *http.Response) (int, []byte, error) {
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, nil, err
	}

	log.Info().Msgf("response body %s", string(responseBody))

	return response.StatusCode, responseBody, nil
}
