package que

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/rs/zerolog/log"
)

type SQSInterface interface {
	GetQueueURL(ctx context.Context, queueName string) (string, error)
	ReceiveMessage(ctx context.Context, sqsUrl string) (*sqs.ReceiveMessageOutput, error)
	SendSQSMessage(ctx context.Context, queueUrl, messageBody string) error
	DeleteMessage(ctx context.Context, queueURL, queMessage string) error
	Subscribe(ctx context.Context, queueURL string, cancel <-chan os.Signal)
	SendMessage(queueName string, messageBody string, chatID, botToken string) error
	SendMessageToUploadWoltImagesToS3(queueName string, menuId string) error
	SendSQSMessageToFIFO(ctx context.Context, queueUrl, messageBody, deduplicationID string) error
}

type SQS struct {
	cli *sqs.Client
}

func NewSQS(cli *sqs.Client) SQSInterface {
	return &SQS{
		cli: cli,
	}
}

func (s SQS) GetQueueURL(ctx context.Context, queueName string) (string, error) {

	result, err := s.cli.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return "", err
	}
	return *result.QueueUrl, nil
}

func (s SQS) ReceiveMessage(ctx context.Context, sqsUrl string) (*sqs.ReceiveMessageOutput, error) {

	msgResult, err := s.cli.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		QueueUrl:            &sqsUrl,
		MaxNumberOfMessages: 1,
	})
	if err != nil {
		return nil, err
	}

	log.Info().Msg(fmt.Sprintf("sqs URL - %s, result sqs - %+v", sqsUrl, msgResult))

	if msgResult == nil || len(msgResult.Messages) == 0 {
		return nil, errors.New("no message found")
	}

	return msgResult, nil
}

func (s SQS) SendSQSMessage(ctx context.Context, queueUrl, messageBody string) error {

	_, err := s.cli.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &queueUrl,
		MessageBody: aws.String(messageBody),
	})

	return err
}

func (s SQS) DeleteMessage(ctx context.Context, queueURL, queMessage string) error {

	_, err := s.cli.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &queueURL,
		ReceiptHandle: &queMessage,
	})
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("could not delete queue message from url - %v", queueURL))
		return errors.New("could not delete queue message")
	}

	return nil
}

func (s SQS) Subscribe(ctx context.Context, queueURL string, cancel <-chan os.Signal) {
	for {
		messages, err := s.ReceiveMessage(ctx, queueURL)
		if err != nil {
			return
		}
		for _, msg := range messages.Messages {
			if msg.Body == nil {
				continue
			}
			go s.DeleteMessage(ctx, queueURL, *msg.ReceiptHandle)
		}

		select {
		case <-cancel:
			return
		case <-time.After(100 * time.Millisecond):
			return
		}
	}
}

func (s SQS) SendMessageToUploadWoltImagesToS3(queueName string, menuId string) error {
	DATA_TYPE := "String"
	gQInput := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	}

	messageAttributes := make(map[string]types.MessageAttributeValue)

	attributes := types.MessageAttributeValue{
		DataType:    aws.String(DATA_TYPE),
		StringValue: aws.String(menuId),
	}

	messageAttributes["menu_id"] = attributes

	result, err := s.cli.GetQueueUrl(
		context.TODO(),
		gQInput,
	)
	if err != nil {
		return err
	}

	msgInput := &sqs.SendMessageInput{
		MessageBody:       aws.String("trigger cmd/menu/wolt_images_upload_to_s3 lambda"),
		QueueUrl:          result.QueueUrl,
		MessageAttributes: messageAttributes,
	}

	if _, err := s.cli.SendMessage(
		context.TODO(),
		msgInput,
	); err != nil {
		return err
	}

	return nil
}

func (s SQS) SendMessage(queueName string, messageBody string, chatID, telegramBotToken string) error {
	DATA_TYPE := "String"
	gQInput := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	}

	messageAttributes := make(map[string]types.MessageAttributeValue)

	attributes := types.MessageAttributeValue{
		DataType:    aws.String(DATA_TYPE),
		StringValue: aws.String(chatID),
	}

	messageAttributes["chat_id"] = attributes

	if telegramBotToken != "" {
		telegramBotAttribute := types.MessageAttributeValue{
			DataType:    aws.String(DATA_TYPE),
			StringValue: aws.String(telegramBotToken),
		}
		messageAttributes["telegram_bot_token"] = telegramBotAttribute
	}

	result, err := s.cli.GetQueueUrl(
		context.TODO(),
		gQInput,
	)
	if err != nil {
		return err
	}

	msgInput := &sqs.SendMessageInput{
		MessageBody:       aws.String(messageBody),
		QueueUrl:          result.QueueUrl,
		MessageAttributes: messageAttributes,
	}

	if _, err := s.cli.SendMessage(
		context.TODO(),
		msgInput,
	); err != nil {
		return err
	}

	return nil
}

func (s SQS) SendSQSMessageToFIFO(ctx context.Context, queueUrl, messageBody, deduplicationID string) error {

	_, err := s.cli.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:               &queueUrl,
		MessageBody:            aws.String(messageBody),
		MessageGroupId:         aws.String(deduplicationID),
		MessageDeduplicationId: aws.String(deduplicationID),
	})

	return err
}
