package general

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Configuration struct {
	AwsConfiguration
	DataStoreConfiguration
	RetryConfiguration
	Sentry

	GlovoConfiguration
	WoltConfiguration
	IIKOConfiguration
	Express24Configuration
	DeliverooConfiguration
	RkeeperConfiguration
	RKeeper7XMLConfiguration
	YarosConfiguration
	PosterConfiguration
	TalabatConfiguration
	PosistConfiguration
	QueConfiguration
	NotificationConfiguration
	FirebaseConfiguration
	WhatsAppConfiguration
	PrometheusCfg
	LegalEntityPaymentCfg
	EmenuConfiguration
	WhatsappBusinessConfiguration
	GoogleSheetsConfiguration
	SmsConfiguration
	StarterAppConfiguration
	IntegrationBaseURL     string `json:"integration_api_base_url"`
	AdminBaseURL           string `json:"admin_api_base_url"`
	SecretEnvironments     string `env:"Prod_Env" envDefault:"ProdEnvs"`
	Stage                  string `json:"stage" envDefault:"dev"`
	TimeZone               string `json:"tz"`
	AppSecret              string `json:"secret"`
	AwsSession             *session.Session
	KwaakaAdminToken       string `json:"kwaaka_admin_token"`
	KwaakaQrMenuToken      string `json:"kwaaka_qr_menu_token"`
	KwaakaWppBusinessToken string `json:"kwaaka_wpp_business_token"`
	KwaakaFilesBucket      string `json:"s3_kwaaka_files_bucket"`
	KwaakaFilesBaseUrl     string `json:"kwaaka_files_base_url"`
	ShaurmaFoodToken       string `json:"shaurma_food_token"`
	GourmetToken           string `json:"gourmet_token"`
}
type StarterAppConfiguration struct {
	BaseUrl string `json:"starter_app_base_url"`
	Token   string `json:"starter_app_token"`
}

type PosistConfiguration struct {
	BaseUrl string `json:"posist_base_url"`
}

type Sentry struct {
	IntegrationDSN string `json:"integration_dsn"`
}

type MulticardConfiguration struct {
	ApplicationId string `json:"multicard_application_id"`
	Secret        string `json:"multicard_secret"`
	StoreId       string `json:"multicard_store_id"`
}

type WhatsappBusinessConfiguration struct {
	BaseUrl    string `json:"wpp_business_base_url"`
	BusinessId string `json:"wpp_business_phone_id"`
	Token      string `json:"wpp_business_token"`
}

type WhatsAppConfiguration struct {
	Instance  string `json:"whatsapp_instance"`
	AuthToken string `json:"whatsapp_auth_token"`
	BaseUrl   string `json:"whatsapp_base_url"`
}

type RkeeperConfiguration struct {
	BaseURL string `json:"rkeeper_base_url"`
	ApiKey  string `json:"rkeeper_api_key"`
}

type TalabatConfiguration struct {
	MiddlewareBaseURL string `json:"talabat_middleware_base_url"`
	MenuBaseUrl       string `json:"talabat_menu_base_url"`
}

type RetryConfiguration struct {
	Count     RetryCount `json:"max_retry_count"`
	QueueName string     `json:"queue_retry"`
}

type RetryCount string

func (rc RetryCount) ToInt() (int, error) {
	count, err := strconv.Atoi(string(rc))
	if err != nil {
		return 0, err
	}

	return count, nil
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

type DeliverooConfiguration struct {
	BaseURL string `json:"deliveroo_base_url"`
	Token   string `json:"deliveroo_token"`
}

type IIKOConfiguration struct {
	BaseURL                 string `json:"iiko_cloud_base_url"`
	TransportToFrontTimeout string `json:"iiko_transport_to_front_timeout"`
}

type RKeeper7XMLConfiguration struct {
	LicenseBaseURL string `json:"rkeeper7_xml_license_base_url"`
}

type PosterConfiguration struct {
	BaseURL string `json:"poster_base_url"`
}

type YarosConfiguration struct {
	BaseURL string `json:"yaros_base_url"`
	Token   string `json:"yaros_token"`
}

type BurgerKingConfiguration struct {
	BaseURL string `json:"burger_king_base_url"`
}

type Express24Configuration struct {
	Token   string `json:"express_24_token"`
	BaseURL string `json:"express24_base_url"`
}

type NotificationConfiguration struct {
	TelegramChatID                            string `json:"notification_bot_chat_id"`
	TelegramChatToken                         string `json:"notification_bot_token"`
	YandexErrorChatID                         string `json:"yandex_error_chat_id"`
	OrderBotToken                             string `json:"order_bot_token"`
	KwaakaDirectTelegramChatId                string `json:"kwaaka_direct_telegram_chat_id"`
	KwaakaDirectNoCourierTelegramChatID       string `json:"kwaaka_direct_no_courier_telegram_chat_id"`
	KwaakaDirectRefundChatID                  string `json:"kwaaka_direct_refund_chat_id"`
	KwaakaDirectKwaakaAdminCompensationChatID string `json:"kwaaka_direct_kwaaka_admin_compensation_chat_id"`
	KwaakaDirect3plNotificationsChatID        string `json:"kwaaka_direct_3pl_notifications_chat_id"`
	AutoUpdatePublicateNotificationChatID     string `json:"auto_update_publicate_notification_chat_id"`
	PutProductToStopListWithErrSolutionChatID string `json:"put_product_to_stop_list_with_err_solution_chat_id"`
	OrderStatChatID                           string `json:"order_stat_chat_id"`
}

type FirebaseConfiguration struct {
	S3BucketName string `json:"firebase_s3_bucket_name"`
	S3FileKey    string `json:"firebase_s3_file_key"`
}
type QueConfiguration struct {
	Telegram      string `json:"queue_telegram"`
	OfflineOrders string `json:"offline_orders"`
}

type PrometheusCfg struct {
	URL      string `json:"prometheus_url"`
	JobName  string `json:"prometheus_job_name"`
	Username string `json:"prometheus_username"`
	Password string `json:"prometheus_password"`
}

type EmenuConfiguration struct {
	EmenuAuthToken                string `json:"emenu_auth_token"`
	EmenuWebhookProductStoplist   string `json:"emenu_webhook_product_stoplist"`
	EmenuWebhookAttributeStoplist string `json:"emenu_webhook_attribute_stoplist"`
	EmenuWebhookURL               string `json:"emenu_webhook_url"`
}

const (
	SECRET_ENV = "SECRET_ENV"
	REGION     = "REGION"
	SENTRY     = "SENTRY"
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

	sentryValue, err := secretClient.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(os.Getenv(SENTRY)),
	})
	if err != nil {
		return Configuration{}, err
	}

	var sentry Sentry

	if err = json.NewDecoder(strings.NewReader(*sentryValue.SecretString)).Decode(&sentry); err != nil {
		return Configuration{}, err
	}

	sqsClient := sqs.NewFromConfig(opts.AwsConfig)

	opts.SqsClient = sqsClient
	opts.Sentry = sentry

	opts.AwsSession = session.Must(session.NewSession())

	return opts, nil
}

type LegalEntityPaymentCfg struct {
	DBCfg LegalEntityPaymentDB
}

type LegalEntityPaymentDB struct {
	Host     string `json:"legal_entity_payment_db_host"`
	Port     string `json:"legal_entity_payment_db_port"`
	User     string `json:"legal_entity_payment_db_user"`
	Password string `json:"legal_entity_payment_db_password"`
	DBName   string `json:"legal_entity_payment_db_name"`
}

type GoogleSheetsConfiguration struct {
	GoogleSheetsEmail     string `json:"google_sheets_email"`
	GoogleSheetsKey       string `json:"google_sheets_key"`
	GoogleSheetsCredsType string `json:"google_sheets_credentials_type"`
}

type SmsConfiguration struct {
	SmsLogin    string `json:"sms_service_login"`
	SmsPassword string `json:"sms_service_password"`
}
