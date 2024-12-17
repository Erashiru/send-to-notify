package config

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Configuration struct {
	AwsConfiguration
	DataStoreConfiguration

	IIKOConfiguration
	RetryConfiguration
	NotificationConfiguration
	QueueConfiguration

	PrometheusCfg

	SecretEnvironments string `env:"Prod_Env" envDefault:"ProdEnvs"`
	Stage              string `json:"stage" envDefault:"dev"`
	TimeZone           string `json:"tz"`
	AdminBaseURL       string `json:"admin_api_base_url"`
	IntegrationBaseURL string `json:"integration_api_base_url"`
	AuthToken          string `json:"admin_api_auth_token"`
}

type RetryConfiguration struct {
	Count     string `json:"max_retry_count"`
	QueueName string `json:"queue_retry"`
}

func (c *RetryConfiguration) ToInt() (int, error) {
	count, err := strconv.Atoi(c.Count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

type QueueConfiguration struct {
	QueueName string `json:"queue_telegram"`
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

type IIKOConfiguration struct {
	BaseURL                 string `json:"iiko_cloud_base_url"`
	TransportToFrontTimeout string `json:"iiko_transport_to_front_timeout"`
	Token                   string `json:"token"`
}

type BurgerKingConfiguration struct {
	BaseURL string `json:"burger_king_base_url"`
}

type NotificationConfiguration struct {
	TelegramChatID    string `json:"notification_bot_chat_id"`
	TelegramChatToken string `json:"notification_bot_token"`
}

type PrometheusCfg struct {
	URL      string `json:"prometheus_url"`
	JobName  string `json:"prometheus_job_name"`
	Username string `json:"prometheus_username"`
	Password string `json:"prometheus_password"`
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

	if err = json.NewDecoder(strings.NewReader(*secretValue.SecretString)).Decode(&opts); err != nil {
		return Configuration{}, err
	}

	return opts, nil
}
