package conf

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type RedisConfiguration struct {
	RedisAddress string `mapstructure:"REDIS_ADDRESS"`
	RedisTtl     string `mapstructure:"REDIS_TTL"`
}

type DatabaseConfiguration struct {
	DbUser     string `mapstructure:"DB_USER"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
	DbHost     string `mapstructure:"DB_HOST"`
	DbPort     string `mapstructure:"DB_PORT"`
	DbName     string `mapstructure:"DB_NAME"`
}

type Configuration struct {
	Database DatabaseConfiguration `mapstructure:",squash"`
	Redis    RedisConfiguration    `mapstructure:",squash"`
	AppName  string                `mapstructure:"APP_NAME"`
	AppPort  uint32                `mapstructure:"APP_PORT"`
	Env      string                `mapstructure:"ENV"`
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

	fmt.Println("DB url", configuration.Database.DbHost)
}

func GetConfiguration() *Configuration {
	return &configuration
}

func GetRedisConnectionURL() string {
	return configuration.Redis.RedisAddress
}
