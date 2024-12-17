package config

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"os"
	"strings"
)

const (
	SECRET_ENV = "SECRET_ENV"
	REGION     = "REGION"
	PROTOCOL   = "http"
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

type Configuration struct {
	AwsConfiguration
	DataStoreConfiguration

	GlovoConfiguration
	WoltConfiguration
	PalomaConfiguration
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
	BaseURL  string `json:"wolt_base_url"`
	Username string `json:"-"`
	Password string `json:"-"`
	ApiKey   string `json:"-"`
}

type PalomaConfiguration struct {
	BaseURL string `json:"paloma_base_url"`
	Class   string `json:"paloma_class"`
	ApiKey  string `json:"paloma_api_key"`
}
