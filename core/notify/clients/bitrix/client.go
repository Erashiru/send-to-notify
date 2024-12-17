package bitrix

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/notify/config"
	models2 "github.com/kwaaka-team/orders-core/core/notify/models"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

const (
	retriesNumber   = 5
	retriesWaitTime = 1 * time.Second
)

const (
	acceptHeader      = "Accept"
	contentTypeHeader = "Content-Type"

	jsonType = "application/json"
)

type Bitrix interface {
	GetTaskItemList(ctx context.Context) (models2.TasksListResponse, error)
	CreateTask(ctx context.Context, title, description string) (models2.AddTaskResponse, error)
	GetLeadByID(ctx context.Context, leadID string) (models2.GetLeadResponse, error)
}

type bitrix struct {
	baseUrl       string
	userID        string
	secret        string
	groupID       string
	responsibleID string
	duties        string

	cli *resty.Client
}

func NewBitrix(cfg config.BitrixConfiguration) Bitrix {

	restyCli := resty.New().
		SetBaseURL(cfg.BaseURL).
		SetRetryCount(retriesNumber).
		SetRetryWaitTime(retriesWaitTime).
		SetHeaders(map[string]string{
			contentTypeHeader: jsonType,
			acceptHeader:      jsonType,
		})

	return &bitrix{
		baseUrl:       cfg.BaseURL,
		userID:        cfg.UserID,
		secret:        cfg.Secret,
		groupID:       cfg.GroupID,
		responsibleID: cfg.ResponsibleID,
		duties:        cfg.Duties,
		cli:           restyCli,
	}

}

func (b bitrix) GetTaskItemList(ctx context.Context) (models2.TasksListResponse, error) {

	path := fmt.Sprintf("/rest/%s/%s/task.item.getlist", b.userID, b.secret)

	var (
		res     models2.TasksListResponse
		errResp models2.ErrorResponse
	)

	rsp, err := b.cli.R().
		SetContext(ctx).
		SetHeader(contentTypeHeader, jsonType).
		SetError(&errResp).
		SetResult(&res).
		Post(path)

	if err != nil {
		return models2.TasksListResponse{}, errResp
	}

	if rsp.IsError() {
		return models2.TasksListResponse{}, errResp
	}

	return res, nil
}

func (b bitrix) CreateTask(ctx context.Context, title, description string) (models2.AddTaskResponse, error) {

	path := fmt.Sprintf("/rest/%s/%s/task.item.add.json", b.userID, b.secret)

	var (
		req     []models2.AddTask
		res     models2.AddTaskResponse
		errResp models2.ErrorResponse
	)

	deadline := time.Now().Add(time.Hour * 1)

	responsibleId, err := strconv.Atoi(b.getDuty(ctx))
	if err != nil {
		log.Err(err).Msgf("%s", b.responsibleID)
	}

	req = append(req, models2.AddTask{
		Title:         title,
		Description:   description,
		GroupID:       b.groupID,
		ResponsibleID: responsibleId,
		Deadline:      &deadline,
	})

	rsp, err := b.cli.R().
		SetContext(ctx).
		SetBody(req).
		SetError(&errResp).
		SetResult(&res).
		Post(path)

	if err != nil {
		return models2.AddTaskResponse{}, errResp
	}

	if rsp.IsError() {
		return models2.AddTaskResponse{}, errResp
	}

	return res, nil
}

func (b bitrix) getDuty(ctx context.Context) string {

	duties := make(map[string]string)

	if err := json.NewDecoder(strings.NewReader(b.duties)).Decode(&duties); err != nil {
		return b.responsibleID
	}

	day := time.Now().Format("2006-01-02")

	responsibleId := b.responsibleID
	if val, ok := duties[day]; ok {
		responsibleId = val
	}

	return responsibleId
}

func (b bitrix) GetLeadByID(ctx context.Context, leadID string) (models2.GetLeadResponse, error) {
	path := fmt.Sprintf("/rest/%s/%s/crm.lead.get.json?id=%s", b.userID, b.secret, leadID)

	var (
		result      models2.GetLeadResponse
		errResponse models2.ErrorResponse
	)

	resp, err := b.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&result).Get(path)
	if err != nil {
		return models2.GetLeadResponse{}, err
	}

	log.Info().Msgf("request bitrix service get lead by id %+v\n", resp.Request.URL)

	if resp.IsError() {
		return models2.GetLeadResponse{}, fmt.Errorf("error get bitrix lead by id %v\n", resp.Error())
	}

	log.Info().Msgf("response bitrix service get lead by id %+v\n", string(resp.Body()))

	return result, nil
}
