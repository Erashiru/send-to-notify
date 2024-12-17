package bitrix

import (
	"github.com/kwaaka-team/orders-core/core/notify/clients/bitrix"
	"github.com/kwaaka-team/orders-core/core/notify/config"
	"github.com/kwaaka-team/orders-core/service/whatsapp"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

type Service interface {
	SendMessageToCustomerFromBitrixLead(ctx context.Context, leadID string) error
}

type ServiceImpl struct {
	bitrixCliDonerSatbayev bitrix.Bitrix
	whatsAppService        whatsapp.Service
}

func NewBitrixService(whatsAppService whatsapp.Service) *ServiceImpl {
	bitrixCliDonerSatbayev := bitrix.NewBitrix(config.BitrixConfiguration{
		BaseURL: "https://dns.bitrix24.kz",
		UserID:  "1235",
		Secret:  "uasben6649pgfnu2",
	})
	return &ServiceImpl{
		bitrixCliDonerSatbayev: bitrixCliDonerSatbayev,
		whatsAppService:        whatsAppService,
	}
}

func (s *ServiceImpl) SendMessageToCustomerFromBitrixLead(ctx context.Context, leadID string) error {
	log.Info().Msgf("send message from bitrix customer, lead id: %s", leadID)

	lead, err := s.bitrixCliDonerSatbayev.GetLeadByID(ctx, leadID)
	if err != nil {
		log.Err(err).Msgf("failed to get bitrix lead id: %s", leadID)
		return err
	}

	var phone string
	if lead.Result.Phone != nil && len(lead.Result.Phone) != 0 {
		phone = lead.Result.Phone[0].Value
	}

	log.Info().Msgf("get phone form bitrix lead succes phone: %s", phone)

	message := "Здравствуйте! 👋🏻\n\nЕсли вы хотите сделать заказ - пожалуйста, воспользуйтесь этой ссылкой: https://satpai.food/"

	if err := s.whatsAppService.SendMessageFromBaseEnvs(ctx, phone, message); err != nil {
		log.Err(err).Msgf("failed to send message whatsapp in phone: %s", phone)
		return err
	}

	log.Info().Msgf("send sessage succes in phone: %s", phone)

	return nil
}
