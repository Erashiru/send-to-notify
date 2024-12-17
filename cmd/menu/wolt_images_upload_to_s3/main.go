package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/kwaaka-team/orders-core/config/general"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"
	"github.com/kwaaka-team/orders-core/service/aws_s3"
	menuServicePkg "github.com/kwaaka-team/orders-core/service/menu"
	"github.com/rs/zerolog/log"
)

var (
	queueName = "small_wolt_images"
)

func run() error {
	ctx := context.Background()

	cfg, err := general.LoadConfig(ctx)
	if err != nil {
		return err
	}

	db, err := cmd.CreateMongo(ctx, cfg.DSURL, cfg.DSDB)
	if err != nil {
		return err
	}

	session := cmd.GetSession()

	s3Service := aws_s3.NewS3Service(session)

	sqsCli := notifyQueue.NewSQS(sqs.NewFromConfig(cfg.AwsConfig))

	menuRepo, err := menuServicePkg.NewMenuMongoRepository(db)
	if err != nil {
		return err
	}

	menuService, err := menuServicePkg.NewMenuService(menuRepo, nil, s3Service, sqsCli)
	if err != nil {
		return err
	}

	menuId, err := receiveMessage(cfg.SqsClient)
	if err != nil {
		return err
	}

	log.Info().Msgf("menu_id = %s", menuId)

	if err = menuService.UploadImagesInWoltFormat(ctx, menuId); err != nil {
		return err
	}

	return nil
}

func receiveMessage(sqsClient *sqs.Client) (string, error) {
	gQInput := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	}

	urlResult, err := sqsClient.GetQueueUrl(
		context.TODO(),
		gQInput,
	)
	if err != nil {
		return "", err
	}
	gMInput := &sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		QueueUrl:            urlResult.QueueUrl,
		MaxNumberOfMessages: 10,
		VisibilityTimeout:   int32(20),
		WaitTimeSeconds:     int32(20),
	}

	msgResult, err := sqsClient.ReceiveMessage(context.TODO(), gMInput)
	if err != nil {
		return "", err
	}

	if msgResult != nil {
		for _, msg := range msgResult.Messages {
			log.Info().Msgf("read message: %v", msg)

			menuId := *msg.MessageAttributes["menu_id"].StringValue

			_, err = sqsClient.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
				QueueUrl:      urlResult.QueueUrl,
				ReceiptHandle: msg.ReceiptHandle,
			})
			if err != nil {
				log.Err(err)
				continue
			}

			return menuId, nil
		}
	} else {
		log.Info().Msg("No messages found")
	}

	return "", err
}

func main() {
	log.Info().Msg("Starting wolt images upload to S3")

	if cmd.IsLambda() {
		lambda.Start(run)
	} else {
		if err := run(); err != nil {
			return
		}
	}
}
