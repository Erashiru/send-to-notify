package menu

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

	GlovoConfiguration
	WoltConfiguration
	ChocofoodConfiguration
	DeliverooConfiguration
	IIKOConfiguration
	SyrveConfiguration
	RkeeperConfiguration
	RKeeper7XMLConfiguration
	PalomaConfiguration
	JowiConfiguration
	BurgerKingConfiguration
	MoySkladConfiguration
	QueConfiguration
	Express24Configuration
	NotificationConfiguration
	S3_BUCKET
	EMENUConfiguration
	TalabatConfiguration
	YandexConfiguration
	YarosConfiguration
	TillypadConfiguration
	Ytimes
	StarterApp
}

type StarterApp struct {
	BaseUrl string `json:"starter_app_base_url"`
}
type Ytimes struct {
	BaseUrl string `json:"ytimes_base_url"`
	Token   string `json:"ytimes_token"`
}

type YarosConfiguration struct {
	BaseURL   string `json:"yaros_base_url"`
	AuthToken string `json:"yaros_token"`
}

type YandexConfiguration struct {
	BaseURL      string `json:"yandex_base_url"`
	ClientID     string `json:"kwaaka_yandex_client_id"`
	ClientSecret string `json:"kwaaka_yandex_client_secret"`
}

type EMENUConfiguration struct {
	AuthToken string `json:"auth_token"`
}

type RKeeper7XMLConfiguration struct {
	LicenseBaseURL string `json:"rkeeper7_xml_license_base_url"`
}

type QueConfiguration struct {
	Telegram string `json:"queue_telegram"`
}

type AwsConfiguration struct {
	AwsConfig        aws.Config
	SqsClient        *sqs.Client
	Region           string `json:"region_name" env:"AWS_REGION" envDefault:"eu-west-1"`
	AWSCognitoPoolID string `json:"aws_cognito_pool_id"`
}

type DataStoreConfiguration struct {
	DSName string `json:"db_engine" env:"db_engine" envDefault:"mongo"`
	DSDB   string `json:"db_name" env:"db_name" envDefault:"kwaaka"`
	DSURL  string `json:"db_url" env:"db_url"  envDefault:"mongodb://localhost:27017"`
}

type TalabatConfiguration struct {
	MenuBaseURL       string `json:"talabat_menu_base_url"`
	MiddlewareBaseURL string `json:"talabat_middleware_base_url"`
}

type PalomaConfiguration struct {
	BaseURL string `json:"paloma_base_url"`
	Class   string `json:"paloma_class"`
	ApiKey  string `json:"paloma_api_key"`
}

type JowiConfiguration struct {
	ApiSecret string `json:"jowi_api_secret"`
	BaseURL   string `json:"jowi_base_url"`
}

type GlovoConfiguration struct {
	BaseURL string `json:"glovo_base_url"`
	Token   string `json:"glovo_token"`
}

type MoySkladConfiguration struct {
	BaseURL     string `json:"moysklad_base_url"`
	Username    string `json:"-"`
	Password    string `json:"-"`
	ApiKey      string `json:"moysklad_api_key"`
	ProductHref string `json:"moysklad_product_href"`
	ProductType string `json:"moysklad_product_type"`
	Protocol    string `json:"moysklad_protocol"`
	Quantity    string `json:"moysklad_quantity"`
}

type Express24Configuration struct {
	BaseURL  string `json:"express24_base_url"`
	Username string `json:"-"`
	Password string `json:"-"`
	Token    string `json:"-"`
}

type DeliverooConfiguration struct {
	BaseURL  string `json:"deliveroo_base_url"`
	Username string `json:"-"`
	Password string `json:"-"`
}

type WoltConfiguration struct {
	BaseURL  string `json:"wolt_base_url"`
	Username string `json:"-"`
	Password string `json:"-"`
	ApiKey   string `json:"-"`
}

type ChocofoodConfiguration struct {
	BaseURL  string `json:"chocofood_base_url"`
	Username string `json:"chocofood_username"`
	Password string `json:"chocofood_password"`
}

type IIKOConfiguration struct {
	BaseURL                 string `json:"iiko_cloud_base_url"`
	TransportToFrontTimeout string `json:"iiko_transport_to_front_timeout"`
	SuperAdminToken         string `json:"super_admin_token"`
}

type SyrveConfiguration struct {
	BaseURL string `json:"syrve_base_url"`
}

type RkeeperConfiguration struct {
	ApiKey  string `json:"rkeeper_api_key"`
	BaseURL string `json:"rkeeper_base_url"`
}

type BurgerKingConfiguration struct {
	BaseURL string `json:"burger_king_base_url"`
}

type NotificationConfiguration struct {
	TelegramChatID                string `json:"notification_bot_chat_id"`
	TelegramChatToken             string `json:"notification_bot_token"`
	AutoUpdatePriceTelegramChatId string `json:"auto_update_price_telegram_chat_id"`
}

type S3_BUCKET struct {
	KwaakaMenuFilesBucket string `json:"s3_kwaaka_menu_files_bucket"`
	ShareMenuBaseUrl      string `json:"share_menu_base_url"`
}

type TillypadConfiguration struct {
	BaseUrl string `json:"tillypad_base_url"`
}

const (
	SECRET_ENV = "SECRET_ENV"
	REGION     = "REGION"
)

func LoadConfig(ctx context.Context, secretEnv, region string) (Configuration, error) {

	var (
		opts Configuration
		err  error
	)

	if secretEnv == "" {
		secretEnv = os.Getenv(SECRET_ENV)
	}

	if region == "" {
		region = os.Getenv(REGION)
	}

	opts.AwsConfig, err = config.LoadDefaultConfig(ctx,
		config.WithRegion(region))
	if err != nil {
		return Configuration{}, err
	}

	secretClient := secretsmanager.NewFromConfig(opts.AwsConfig)
	secretValue, err := secretClient.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretEnv),
	})
	if err != nil {
		return Configuration{}, err
	}

	if err := json.NewDecoder(strings.NewReader(*secretValue.SecretString)).Decode(&opts); err != nil {
		return Configuration{}, err
	}

	return opts, nil
}
