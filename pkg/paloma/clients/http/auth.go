package http

import (
	"context"
	"fmt"
)

func (p *paloma) Auth(ctx context.Context) error {
	if p.apiKey == "" {
		return ErrApiKey
	}

	if p.class == "" {
		return ErrClass
	}

	p.cli.QueryParam.Set("authkey", p.apiKey)
	p.cli.QueryParam.Set("class", p.class)
	p.cli.SetBaseURL(p.cli.BaseURL + "/company/api")

	fmt.Println(p.cli.BaseURL)

	return nil
}
