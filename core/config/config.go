package config

import (
	"context"
	"encoding/json"
	"github.com/kwaaka-team/orders-core/core/models"
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

	DeliverooConfiguration
	QueConfiguration
	GlovoConfiguration
	WoltConfiguration
	IIKOConfiguration
	SyrveConfiguration
	JowiConfiguration
	PalomaConfiguration
	PosterConfiguration
	BurgerKingConfiguration
	RKeeperConfiguration
	ChocofoodConfiguration
	NotificationConfiguration
	RetryConfiguration
	TalabatConfiguration
	YandexConfiguration
	RKeeper7XMLConfiguration
	YarosConfiguration
	Express24Configuration
	TillypadConfiguration
	WhatsAppConfiguration
	PostgreSqlConfiguration
	IokaConfiguration
	PaymeConfiguration
	WoopPayConfiguration
	WhatsappBusinessConfiguration
	RedisConfig
	QueueUrls
	Kwaaka3pl
	Ytimes
	PosistConfiguration
	KaspiSaleScoutConfiguration
	MulticardConfiguration
	StarterAppConfiguration
	AdminApi
}

type StarterAppConfiguration struct {
	BaseUrl                       string `json:"starter_app_base_url"`
	StarterAppSaleScoutProxyToken string `json:"starter_app_sale_scout_proxy_token"`
	Token                         string `json:"starter_app_token"`
}

type MulticardConfiguration struct {
	ApplicationId string `json:"multicard_application_id"`
	Secret        string `json:"multicard_secret"`
	StoreId       string `json:"multicard_store_id"`
	BaseUrl       string `json:"multicard_base_url"`
}

type RedisConfig struct {
	Addr     string `json:"redis_addr"`
	Username string `json:"redis_username"`
	Password string `json:"redis_password"`
}

type PosistConfiguration struct {
	BaseUrl string `json:"posist_base_url"`
}

type Ytimes struct {
	BaseUrl string `json:"ytimes_base_url"`
	Token   string `json:"ytimes_token"`
}

type Kwaaka3pl struct {
	Kwaaka3plBaseUrl   string `json:"dispatcher_base_url"`
	Kwaaka3plAuthToken string `json:"3pl_api_key"`
	Kwaaka3plQueue     string `json:"3pl_queue_url"`
}

type QueueUrls struct {
	PaymentsQueueUrl         string `json:"payments_queue_url"`
	WhatsappMessagesQueueUrl string `json:"whatsapp_messages_queue_url"`
}
type IokaConfiguration struct {
	BaseUrl string `json:"ioka_base_url"`
	ApiKey  string `json:"ioka_api_key"`
}

type KaspiSaleScoutConfiguration struct {
	BaseUrl    string `json:"kaspi_salescout_base_url"`
	Token      string `json:"kaspi_salescout_api_key"`
	MerchantID string `json:"kaspi_salescout_merchant_id"`
}

type PaymeConfiguration struct {
	BaseUrl string `json:"payme_base_url"`
	ApiKey  string `json:"payme_api_key"`
}

type WoopPayConfiguration struct {
	BaseUrl   string `json:"wooppay_base_url"`
	ResultUrl string `json:"wooppay_result_url"`
}

type TillypadConfiguration struct {
	BaseUrl string `json:"tillypad_base_url"`
}

type SyrveConfiguration struct {
	BaseURL string `json:"syrve_base_url"`
}

type PostgreSqlConfiguration struct {
	ConnectionString string `json:"postgres_stage_connection_string"`
}

type PalomaConfiguration struct {
	BaseURL string `json:"paloma_base_url"`
	Class   string `json:"paloma_class"`
	ApiKey  string `json:"paloma_api_key"`
}

type PosterConfiguration struct {
	BaseURL           string `json:"poster_base_url"`
	ApplicationID     string `json:"poster_application_id"`
	ApplicationSecret string `json:"poster_application_secret"`
	RedirectURI       string `json:"poster_redirect_uri"`
}

type YarosConfiguration struct {
	BaseURL    string `json:"yaros_base_url"`
	InfoSystem string `json:"yaros_info_system"`
	Token      string `json:"yaros_token"`
}

type DeliverooConfiguration struct {
	BaseURL string `json:"deliveroo_base_url"`
}

type RetryConfiguration struct {
	Count     string `json:"max_retry_count"`
	QueueName string `json:"queue_retry"`
}
type QueConfiguration struct {
	Telegram string `json:"queue_telegram"`
}

type RKeeperConfiguration struct {
	RKeeperApiKey  string `json:"rkeeper_api_key"`
	RKeeperBaseURL string `json:"rkeeper_base_url"`
}

type Express24Configuration struct {
	BaseURL string `json:"express24_base_url"`
}

type RKeeper7XMLConfiguration struct {
	LicenseBaseURL string `json:"rkeeper7_xml_license_base_url"`
}

type JowiConfiguration struct {
	ApiKey    string `json:"jowi_api_key"`
	ApiSecret string `json:"jowi_api_secret"`
	BaseURL   string `json:"jowi_base_url"`
}

type AwsConfiguration struct {
	AwsConfig aws.Config
	SqsClient *sqs.Client
	Region    string `json:"region_name" env:"AWS_REGION" envDefault:"eu-west-1"`
}

type DataStoreConfiguration struct {
	DSName string `json:"db_engine" env:"db_engine" envDefault:"mongo"`
	DSDB   string `json:"db_name" env:"db_name" envDefault:"kwaaka"`
	DSURL  string `json:"db_url" env:"db_url"  envDefault:"mongodb://localhost:27017"`
}

type GlovoConfiguration struct {
	BaseURL string `json:"glovo_base_url"`
	Token   string `json:"glovo_token"`
}

type WoltConfiguration struct {
	BaseURL string `json:"wolt_base_url"`
}

type TalabatConfiguration struct {
	MiddlewareBaseURL string `json:"talabat_middleware_base_url"`
	MenuBaseUrl       string `json:"talabat_menu_base_url"`
}

type YandexConfiguration struct {
	BaseURL      string `json:"yandex_base_url"`
	ClientID     string `json:"kwaaka_yandex_client_id"`
	ClientSecret string `json:"kwaaka_yandex_client_secret"`
}

type ChocofoodConfiguration struct {
	BaseURL  string `json:"chocofood_base_url"`
	Username string `json:"chocofood_username"`
	Password string `json:"chocofood_password"`
}

type IIKOConfiguration struct {
	BaseURL                 string `json:"iiko_cloud_base_url"`
	TransportToFrontTimeout string `json:"iiko_transport_to_front_timeout"`
}

type BurgerKingConfiguration struct {
	BaseURL string `json:"burger_king_base_url"`
}

type NotificationConfiguration struct {
	TelegramChatID             string `json:"notification_bot_chat_id"`
	TelegramChatToken          string `json:"notification_bot_token"`
	YandexErrorChatID          string `json:"yandex_error_chat_id"`
	OrderBotToken              string `json:"order_bot_token"`
	KwaakaDirectTelegramChatId string `json:"kwaaka_direct_telegram_chat_id"`
}

type WhatsAppConfiguration struct {
	Instance  string `json:"whatsapp_instance"`
	AuthToken string `json:"whatsapp_auth_token"`
	BaseUrl   string `json:"whatsapp_base_url"`
}

type WhatsappBusinessConfiguration struct {
	BaseUrl    string `json:"wpp_business_base_url"`
	BusinessId string `json:"wpp_business_phone_id"`
	Token      string `json:"wpp_business_token"`
}

type AdminApi struct {
	AdminBaseURL         string `json:"admin_api_base_url"`
	WoltDiscountRunToken string `json:"wolt_discount_run_token"`
}

func LoadConfig(ctx context.Context) (Configuration, error) {

	var opts Configuration
	var err error

	opts.AwsConfig, err = config.LoadDefaultConfig(ctx,
		config.WithRegion(os.Getenv(models.REGION)))
	if err != nil {
		return Configuration{}, err
	}

	secretClient := secretsmanager.NewFromConfig(opts.AwsConfig)
	secretValue, err := secretClient.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(os.Getenv(models.SECRET_ENV)),
	})
	if err != nil {
		return Configuration{}, err
	}

	sqsClient := sqs.NewFromConfig(opts.AwsConfig)

	opts.SqsClient = sqsClient

	if err = json.NewDecoder(strings.NewReader(*secretValue.SecretString)).Decode(&opts); err != nil {
		return Configuration{}, err
	}

	return opts, nil
}
