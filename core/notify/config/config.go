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

const (
	SECRET_ENV = "SECRET_ENV"
	REGION     = "REGION"
)

type Configuration struct {
	AwsConfiguration
	DataStoreConfiguration
	FireBaseConfiguration
	BitrixConfiguration
	ClickUpConfiguration
}

type ClickUpConfiguration struct {
	BaseURL           string `json:"clickup_base_url"`
	OrderErrorsListID string `json:"clickup_order_errors_list"`
	MenuErrorsListID  string `json:"clickup_menu_errors_list"`
	Token             string `json:"clickup_token"`
	Assignee          string `json:"clickup_assignee"`
}

type FireBaseConfiguration struct {
	FireBaseDBURL    string `json:"firebase_db_url" env:"firebase_db_url"`
	FireBaseFilePath string `json:"firebase_secret_file_url"`
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

type BitrixConfiguration struct {
	BaseURL       string `json:"bitrix_base_url"`
	UserID        string `json:"bitrix_user_id"`
	Secret        string `json:"bitrix_token"`
	GroupID       string `json:"bitrix_support_group_id"`
	ResponsibleID string `json:"bitrix_responsible_id"`

	Duties string `json:"bitrix_duties"`
}

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
