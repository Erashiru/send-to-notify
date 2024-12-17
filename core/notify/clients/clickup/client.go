package clickup

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/core/notify/config"
	models2 "github.com/kwaaka-team/orders-core/core/notify/models"
	"strconv"
	"time"
)

const (
	acceptHeader      = "Accept"
	contentTypeHeader = "Content-Type"
	authorization     = "Authorization"
	jsonType          = "application/json"
	bearer            = "Bearer "
)

type ClickUp interface {
	CreateTask(ctx context.Context, message models2.Message) (models2.ClickUpTaskResponse, error)
}

type clickup struct {
	baseUrl           string
	token             string
	orderErrorsListID string
	menuErrorsListID  string
	assignee          string
	cli               *resty.Client
}

func NewClickUp(cfg config.ClickUpConfiguration) ClickUp {

	restyCli := resty.New().
		SetBaseURL(cfg.BaseURL).
		SetHeaders(map[string]string{
			contentTypeHeader: jsonType,
			acceptHeader:      jsonType,
			authorization:     bearer + cfg.Token,
		})

	return &clickup{
		baseUrl:           cfg.BaseURL,
		token:             cfg.Token,
		orderErrorsListID: cfg.OrderErrorsListID,
		menuErrorsListID:  cfg.MenuErrorsListID,
		assignee:          cfg.Assignee,
		cli:               restyCli,
	}
}

func (c clickup) CreateTask(ctx context.Context, message models2.Message) (models2.ClickUpTaskResponse, error) {
	var listID string

	switch message.TaskList {
	case models2.MENU_ERROR:
		listID = c.menuErrorsListID
	case models2.ORDER_ERROR:
		listID = c.orderErrorsListID
	default:
		return models2.ClickUpTaskResponse{}, models2.ErrNoTaskList
	}

	path := fmt.Sprintf("/list/%s/task", listID)

	dueDate := time.Now().Add(time.Hour * 24)
	timestamp := dueDate.Unix()

	var (
		errResp models2.ErrorClickUpTask
		res     models2.ClickUpTaskResponse
	)

	assignee, err := strconv.Atoi(c.assignee)
	if err != nil {
		return models2.ClickUpTaskResponse{}, err
	}

	task := models2.ClickUpTask{
		Name:        message.Title,
		Description: message.Title,
		DueDate:     timestamp,
		Priority:    1,
		Assignees:   []int{assignee},
	}

	rsp, err := c.cli.R().
		SetContext(ctx).
		SetBody(task).
		SetError(&errResp).
		SetResult(&res).
		Post(path)

	if err != nil {
		return models2.ClickUpTaskResponse{}, errResp
	}

	if rsp.IsError() {
		return models2.ClickUpTaskResponse{}, errResp
	}

	return res, nil
}
