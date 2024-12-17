package dto

import (
	models2 "github.com/kwaaka-team/orders-core/core/notify/models"
)

type Message struct {
	Title       string
	Description string
	Services    []Service
	TaskList    TaskList
}

func (m Message) ToModel() models2.Message {

	services := make([]models2.Service, 0, len(m.Services))
	for _, service := range m.Services {
		services = append(services, models2.Service(service))
	}

	return models2.Message{
		Title:       m.Title,
		Description: m.Description,
		Services:    services,
		TaskList:    models2.TaskList(m.TaskList),
	}

}

type Results []Result

type Result struct {
	Service Service
	Status  Status
	Message string
}

func FromResults(results []models2.Result) Results {

	res := make(Results, 0, len(results))

	for _, result := range results {
		res = append(res, Result{
			Status:  Status(result.Status),
			Service: Service(result.Service),
			Message: result.Message,
		})
	}

	return res
}
