package config

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Configuration struct {
	AwsConfiguration
	DataStoreConfiguration

	SecretEnvironments string `env:"Prod_Env" envDefault:"ProdEnvs"`
	Stage              string `json:"stage" envDefault:"dev"`
	TimeZone           string `json:"tz"`
	AppSecret          string `json:"secret"`

	GlovoConfiguration
	WoltConfiguration
	IIKOConfiguration
	JowiConfiguration
}

type AwsConfiguration struct {
	AwsConfig aws.Config
	SqsClient *sqs.Client
	Region    string `json:"region_name" env:"AWS_REGION" envDefault:"eu-west-1"`
}

type DataStoreConfiguration struct {
	// DataStore name (format: mongo/null)
	DSName string `json:"db_engine" env:"db_engine" envDefault:"mongo"`
	// DataStore database name (format: menu)
	DSDB string `json:"db_name" env:"db_name" envDefault:"kwaaka"`
	// DataStore URL (format: mongodb://localhost:27017)
	DSURL string `json:"db_url" env:"db_url"  envDefault:"mongodb://localhost:27017"`
}

type GlovoConfiguration struct {
	BaseURL string `json:"glovo_base_url"`
	Token   string `json:"glovo_token"`
}

type WoltConfiguration struct {
	BaseURL string `json:"wolt_base_url"`
}

type IIKOConfiguration struct {
	BaseURL                 string `json:"iiko_cloud_base_url"`
	TransportToFrontTimeout string `json:"iiko_transport_to_front_timeout"`
}

type JowiConfiguration struct {
	BaseURL string `json:"jowi_base_url"`
}

type BurgerKingConfiguration struct {
	BaseURL string `json:"burger_king_base_url"`
}

type NotificationConfiguration struct {
	TelegramChatID    string `json:"notification_bot_chat_id"`
	TelegramChatToken string `json:"notification_bot_token"`
}

const (
	SECRET_ENV = "SECRET_ENV"
	REGION     = "REGION"
)

func LoadConfig(ctx context.Context) (Configuration, error) {

	var opts Configuration
	var err error

	opts.AwsConfig, err = config.LoadDefaultConfig(ctx,
		config.WithRegion(os.Getenv(REGION)))

	if err != nil {
		return Configuration{}, err
	}

	secretClient := secretsmanager.NewFromConfig(opts.AwsConfig)
	secretValue, err := secretClient.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(os.Getenv(SECRET_ENV)),
	})
	if err != nil {
		return Configuration{}, err
	}

	if err := json.NewDecoder(strings.NewReader(*secretValue.SecretString)).Decode(&opts); err != nil {
		return Configuration{}, err
	}

	return opts, nil
}
