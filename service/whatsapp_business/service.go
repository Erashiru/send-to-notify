package whatsapp_business

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/core/config"
	models2 "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/service/whatsapp_business/models"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"io"
	"net/http"
	"time"
)

var (
	table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
)

const (
	ContentType = "Content-Type"
	JsonType    = "application/json"

	Whatsapp = "whatsapp"

	VerifyTemplate = "verify"

	English = "eu_US"
	Russian = "ru"

	TemplateMsg = "template"
	TextMsg     = "text"

	ButtonWidget = "button"

	IndividualRecipientType = "individual"
)

type Service interface {
	SendVerificationCode(ctx context.Context, phoneNum, restGroupId string) error
}

type ServiceImpl struct {
	wppRestyClient *resty.Client
	redisClient    *redis.Client
	PhoneId        string
}

func NewWppBusinessService(cfg config.WhatsappBusinessConfiguration, redis *redis.Client) (Service, error) {
	wppResty := resty.New().
		SetBaseURL(cfg.BaseUrl).
		SetHeaders(map[string]string{
			ContentType: JsonType,
		}).
		SetAuthToken(cfg.Token)

	return &ServiceImpl{
		wppRestyClient: wppResty,
		redisClient:    redis,
		PhoneId:        cfg.BusinessId,
	}, nil
}

func (s *ServiceImpl) SendVerificationCode(ctx context.Context, phoneNum, restGroupId string) error {
	code, err := generateCode(4)
	if err != nil {
		return err
	}

	req := models2.RedisCodeRestaurantGroup{
		Code:              code,
		RestaurantGroupId: restGroupId,
	}

	data, err := json.Marshal(&req)
	if err != nil {
		return err
	}

	if err := s.redisClient.Set(ctx, phoneNum, data, 5*time.Minute).Err(); err != nil {
		return err
	}

	if err := s.sendCode(ctx, phoneNum, code); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) sendCode(ctx context.Context, phoneNum, code string) error {
	path := fmt.Sprintf("/%s/messages", s.PhoneId)

	body := models.VerificationCodeRequest{
		MessagingProduct: Whatsapp,
		RecipientType:    IndividualRecipientType,
		To:               phoneNum,
		Type:             TemplateMsg,
		Template: models.TemplateData{
			Name: VerifyTemplate,
			Language: models.LanguageData{
				Code: Russian,
			},
			Components: []models.Component{
				{
					Type: "body",
					Parameters: []models.Parameter{
						{Type: TextMsg, Text: code},
					},
				},
				{
					Type:    ButtonWidget,
					SubType: "url",
					Index:   "0",
					Parameters: []models.Parameter{
						{Type: TextMsg, Text: code},
					},
				},
			},
		},
	}

	resp, err := s.wppRestyClient.R().
		SetBody(body).
		SetContext(ctx).
		Post(path)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.New("invalid status code")
	}

	return nil
}

func generateCode(max int) (string, error) {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		return "", err
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b), nil
}
