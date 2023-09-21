package conf

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
)

type AffiliateConfiguration struct {
	RewardProgram    string `mapstructure:"AFF_REWARD_PROGRAM"`
	SellerId         string `mapstructure:"AFF_SELLER_ID"`
	StellaCommission string `mapstructure:"AFF_STELLA_COMMISSION"`
	WithdrawFee      string `mapstructure:"AFF_WITHDRAW_FEE"`
}

type TikiConfiguration struct {
	ApiUrl         string `mapstructure:"TIKI_API_URL"`
	ApiKey         string `mapstructure:"TIKI_API_KEY"`
	ExchangeApiKey string `mapstructure:"TIKI_EXCHANGE_API_KEY"`
}

type RewardShippingConfiguration struct {
	BaseUrl string `mapstructure:"REWARD_SHIPPING_URL"`
	ApiKey  string `mapstructure:"REWARD_SHIPPING_KEY"`
}

type RedisConfiguration struct {
	RedisAddress string `mapstructure:"REDIS_ADDRESS"`
	RedisTtl     string `mapstructure:"REDIS_TTL"`
}

type DatabaseConfiguration struct {
	WriteDbUser     string `mapstructure:"WRITE_DB_USER"`
	ReadDbUser      string `mapstructure:"READ_DB_USER"`
	ReadDbPassword  string `mapstructure:"READ_DB_PASSWORD"`
	WriteDbPassword string `mapstructure:"WRITE_DB_PASSWORD"`
	ReadDbHost      string `mapstructure:"READ_DB_HOST"`
	WriteDbHost     string `mapstructure:"WRITE_DB_HOST"`
	DbPort          string `mapstructure:"DB_PORT"`
	DbName          string `mapstructure:"DB_NAME"`
}

type NotificationConfiguration struct {
	EndPoint               string  `mapstructure:"NOTIFY_ENDPOINT"`
	AccessToken            string  `mapstructure:"NOTIFY_SERVICE_ACCESS_TOKEN"`
	FrontendClaimURL       string  `mapstructure:"FE_CLAIM_URL"`
	NotifyMethod           string  `mapstructure:"NOTIFY_METHOD"`
	SellerBalanceThreshold float64 `mapstructure:"SELLER_BALANCE_THRESHOLD"`
	WorkerBalanceThreshold float64 `mapstructure:"WORKER_BALANCE_THRESHOLD"`
}

type EvmRpcEndpointConfiguration struct {
	EVMChainID  int64  `mapstructure:"EVM_CHAIN_ID"`
	EndPoint    string `mapstructure:"EVMRPC_ENDPOINT"`
	ExplorerURL string `mapstructure:"EXPLORER_URL"`
}

type KafkaConfiguration struct {
	KafkaURL    string `mapstructure:"KAFKA_URL" yaml:"kafkaUrl" toml:"kafkaUrl" xml:"kafkaUrl" json:"kafkaUrl,omitempty"`
	User        string `mapstructure:"KAFKA_USER" yaml:"user" toml:"user" xml:"user" json:"user,omitempty"`
	Password    string `mapstructure:"KAFKA_PASSWORD" yaml:"password" toml:"password" xml:"password" json:"password,omitempty"`
	AuthenType  string `mapstructure:"KAFKA_AUTHEN_TYPE" yaml:"authenType" toml:"authenType" xml:"authenType" json:"authenType,omitempty"`
	CaCertPath  string `mapstructure:"KAFKA_CA_CERT_PATH" yaml:"caCertPath" toml:"caCertPath" xml:"caCertPath" json:"caCertPath,omitempty"`
	TlsCertPath string `mapstructure:"KAFKA_TLS_CERT_PATH" yaml:"tlsCertPath" toml:"tlsCertPath" xml:"tlsCertPath" json:"tlsCertPath,omitempty"`
	TlsKeyPath  string `mapstructure:"KAFKA_TLS_KEY_PATH" yaml:"tlsKeyPath" toml:"tlsKeyPath" xml:"tlsKeyPath" json:"tlsKeyPath,omitempty"`
}

type DiscordConfig struct {
	// Token is the OAuth token for pushing messages.
	Token       string   `mapstructure:"DISCORD_TOKEN" yaml:"token" toml:"token" xml:"token" json:"token,omitempty"`
	SubChannels []string `mapstructure:"DISCORD_SUB_CHANNELS" yaml:"subChannels" toml:"subChannels" xml:"subChannels" json:"subChannels,omitempty"`
}
type WebhookConfiguration struct {
	DcConf DiscordConfig `mapstructure:",squash"`
}

type Configuration struct {
	Database          DatabaseConfiguration       `mapstructure:",squash"`
	Redis             RedisConfiguration          `mapstructure:",squash"`
	Notify            NotificationConfiguration   `mapstructure:",squash"`
	EvmRpc            EvmRpcEndpointConfiguration `mapstructure:",squash"`
	Kafka             KafkaConfiguration          `mapstructure:",squash"`
	Webhook           WebhookConfiguration        `mapstructure:",squash"`
	RewardShipping    RewardShippingConfiguration `mapstructure:",squash"`
	Aff               AffiliateConfiguration      `mapstructure:",squash"`
	Tiki              TikiConfiguration           `mapstructure:",squash"`
	AccessTradeAPIKey string                      `mapstructure:"ACCESSTRADE_APIKEY"`
	CreatorBEToken    string                      `mapstructure:"CREATOR_BE_TOKEN"`
	AppName           string                      `mapstructure:"APP_NAME"`
	AppPort           uint32                      `mapstructure:"APP_PORT"`
	Env               string                      `mapstructure:"ENV"`
	WebhookEndPoint   string                      `mapstructure:"WEBHOOK_ENDPOINT"`
	CreatorAuthUrl    string                      `mapstructure:"CREATOR_AUTH_URL"`
	AppAuthUrl        string                      `mapstructure:"APP_AUTH_URL"`
}

var configuration Configuration

func init() {
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}
	fmt.Println(envFile)
	viper.SetConfigFile("./.env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		viper.SetConfigFile(fmt.Sprintf("../%s", envFile))
		if err := viper.ReadInConfig(); err != nil {
			log.Logger.Printf("Error reading config file \"%s\", %v", envFile, err)
		}
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Logger.Fatal().Msgf("Unable to decode config into map, %v", err)
	}

	fmt.Println("EVM ChainId:", configuration.EvmRpc.EVMChainID)
	fmt.Println("EVM Rpc URL:", configuration.EvmRpc.EndPoint)
	fmt.Println("DB url", configuration.Database.WriteDbHost)
}

func GetConfiguration() *Configuration {
	return &configuration
}

func GetRedisConnectionURL() string {
	return configuration.Redis.RedisAddress
}
