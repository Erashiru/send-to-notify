package v1

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/kwaaka-team/orders-core/core/wolt_discount_run/service"
	"github.com/rs/zerolog/log"
)

type Server struct {
	service *service.Service
}

func NewServer(service *service.Service) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) SqsProxy(ctx context.Context, req events.SQSEvent) error {
	if req.Records != nil {
		log.Printf("Request sqs event : %+v\n", req)
		s.WoltDiscountRun(req)
	}
	return nil
}
