package sms

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	models2 "github.com/kwaaka-team/orders-core/core/models"
	storeGroupServicePkg "github.com/kwaaka-team/orders-core/service/storegroup"
	"github.com/redis/go-redis/v9"
	"io"
	"time"
)

type Service interface {
	SendVerificationCode(ctx context.Context, phoneNum, restGroupId string) (string, string, error)
	SendMessage(ctx context.Context, phoneNumber string, message string) error
	IsSmsServiceError(err error) (bool, int, string)
}

type ServiceImpl struct {
	SmsClient         *resty.Client
	redisClient       *redis.Client
	storeGroupService storeGroupServicePkg.Service
}

const (
	BASEURL = "https://smsc.kz/sys/send.php"
)

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func NewSmsService(SmsLogin, SmsPsw string, redisCli *redis.Client, service storeGroupServicePkg.Service) (Service, error) {

	smsCli := resty.New().
		SetBaseURL(BASEURL).
		SetQueryParams(map[string]string{
			"login": SmsLogin,
			"psw":   SmsPsw,
			"fmt":   "3",
		})

	return &ServiceImpl{SmsClient: smsCli, redisClient: redisCli, storeGroupService: service}, nil
}

func (s *ServiceImpl) SendMessage(ctx context.Context, phoneNumber string, message string) error {
	var errResp Response

	resp, err := s.SmsClient.R().
		SetQueryParams(map[string]string{
			"phones": phoneNumber,
			"mes":    message,
		}).SetContext(ctx).SetResult(&errResp).Post("")
	if err != nil {
		return err
	}

	if resp.IsError() {
		return err
	}

	if err := s.isErrorResponse(&errResp); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) SendVerificationCode(ctx context.Context, phoneNum, restGroupId string) (string, string, error) {
	restGroup, err := s.storeGroupService.GetStoreGroupByID(ctx, restGroupId)
	if err != nil {
		return "", "", err
	}

	code, err := generateCode(4)
	if err != nil {
		return "", "", err
	}

	req := models2.RedisCodeRestaurantGroup{
		Code:              code,
		RestaurantGroupId: restGroupId,
	}

	data, err := json.Marshal(&req)
	if err != nil {
		return "", "", err
	}

	if err = s.redisClient.Set(ctx, phoneNum, data, 5*time.Minute).Err(); err != nil {
		return "", "", err
	}

	if err = s.SendMessage(ctx, phoneNum, fmt.Sprintf("%s: код подтверждения в %s", code, restGroup.Name)); err != nil {
		return code, restGroup.Name, err
	}

	return code, restGroup.Name, nil
}

func (s *ServiceImpl) isErrorResponse(resp *Response) error {
	if resp == nil || resp.Error == "" {
		return nil
	}

	var err Error
	switch resp.ErrorCode {
	case 1:
		err = ParamErr
	case 2:
		err = CredsErr
	case 3:
		err = MoneyErr
	case 4:
		err = IpErr
	case 5:
		err = RestrictionErr
	case 6:
		err = DateFormatErr
	case 7:
		err = InvalidNumErr
	case 8:
		err = SendErr
	case 9:
		err = SpamErr
	default:
		err = Error{Code: resp.ErrorCode, Message: "Unknown error from smsc.kz"}
	}

	return fmt.Errorf("message %d ended with error: %w and with description %s", resp.SmsId, err, resp.Error)
}

func (s *ServiceImpl) IsSmsServiceError(err error) (bool, int, string) {
	var smsErr Error
	if errors.As(err, &smsErr) {
		return true, smsErr.Code, smsErr.Message
	}
	return false, 0, ""
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
