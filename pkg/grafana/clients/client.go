package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	"github.com/kwaaka-team/orders-core/pkg/grafana/config"
	httpService "github.com/kwaaka-team/orders-core/pkg/net-http-client/http"
	"time"
)

type Grafana interface {
	AddAnnotations(ctx context.Context, dashboardUID, description string, panelId int, tags []string) error
}

type grafanaImpl struct {
	client httpService.Client
}

func NewGrafanaClient(config config.Config) *grafanaImpl {
	return &grafanaImpl{
		client: httpService.NewHTTPClient(config.BaseURL),
	}
}

type Request struct {
	DashboardUID string   `json:"dashboardUID"`
	PanelId      int      `json:"panelId"`
	Time         int64    `json:"time"`
	TimeEnd      int64    `json:"timeEnd"`
	Tags         []string `json:"tags"`
	Text         string   `json:"text"`
}

type Response struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}

func (grafana *grafanaImpl) AddAnnotations(ctx context.Context, dashboardUID, description string, panelId int, tags []string) error {
	req := Request{
		DashboardUID: dashboardUID,
		PanelId:      panelId,
		Time:         time.Now().UnixMilli(),
		TimeEnd:      time.Now().UnixMilli(),
		Tags:         tags,
		Text:         description,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	status, response, err := grafana.client.Post("/api/annotations", body, map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		return err
	}

	if status >= 400 {
		return fmt.Errorf("response error, status %d", status)
	}

	var result Response

	if err = json.Unmarshal(response, &result); err != nil {
		return err
	}

	utils.Beautify("response body", result)

	return nil
}
