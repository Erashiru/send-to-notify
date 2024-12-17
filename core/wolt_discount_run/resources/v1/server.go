package v1

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/kwaaka-team/orders-core/core/wolt_discount_run/models"
	"github.com/rs/zerolog/log"
)

func (s *Server) WoltDiscountRun(req events.SQSEvent) {
	for _, message := range req.Records {
		var request models.DiscountRunRequest
		if err := json.Unmarshal([]byte(message.Body), &request); err != nil {
			log.Err(err).Msgf("unmarshalling message error")
			continue
		}

		if err := s.service.WoltDiscountRun(context.Background(), request); err != nil {
			log.Err(err).Msgf("wolt discount run error with: %+v", request)
			continue
		}
	}
}
