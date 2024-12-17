package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/pkg/multicard"
	"github.com/kwaaka-team/orders-core/pkg/multicard/dto"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func main() {
	if cmd.IsLambda() {
		lambda.Start(run)
	} else {
		if err := run(context.Background()); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func run(ctx context.Context) error {
	const op = "mutricard.refreshToken"

	encoderCfg := zap.NewProductionConfig()
	encoderCfg.EncoderConfig.TimeKey = "timestamp"
	encoderCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncoderConfig.StacktraceKey = ""
	l, err := encoderCfg.Build()
	if err != nil {
		return err
	}
	logger := l.Sugar()
	defer logger.Sync()

	log := logger.With(zap.Fields(
		zap.String("op", op),
	))

	opts, err := config.LoadConfig(context.Background())
	if err != nil {
		return err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     opts.RedisConfig.Addr,
		Username: opts.RedisConfig.Username,
		Password: opts.RedisConfig.Password,
	})
	if err = redisClient.Ping(ctx).Err(); err != nil {
		log.Error("failed to connect to redis", zap.Error(err))
		return err
	}

	var (
		url     = "/auth"
		resp    = dto.AuthResponse{}
		cli     = resty.New().SetBaseURL(opts.MulticardConfiguration.BaseUrl)
		authReq = dto.AuthRequest{
			ApplicationId: opts.MulticardConfiguration.ApplicationId,
			Secret:        opts.MulticardConfiguration.Secret,
		}
	)

	r, err := cli.R().
		SetBody(authReq).
		SetContext(ctx).
		SetResult(&resp).
		Post(url)
	if err != nil {
		log.Error("failed to create auth request", zap.Error(err))
		return err
	}

	if r.IsError() {
		log.Error("auth request failed wit fail response", zap.Any("response", resp))
		return err
	}

	if err = redisClient.Set(ctx, multicard.AuthKey, resp.Token, 24*time.Hour).Err(); err != nil {
		return err
	}

	return nil
}
