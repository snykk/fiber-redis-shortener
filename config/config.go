package config

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

var AppConfig Config

type Config struct {
	Host        string
	Port        int
	Environment string
	Debug       bool

	REDISHost     string
	REDISPassword string
	REDISDbno     int
	REDISExpired  int
}

func InitializeAppConfig() error {
	viper.SetConfigName(".env") // allow directly reading from .env file
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/")
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	// assign value
	AppConfig.Host = viper.GetString("HOST")
	AppConfig.Port = viper.GetInt("PORT")
	AppConfig.Environment = viper.GetString("ENVIRONMENT")
	AppConfig.Debug = viper.GetBool("DEBUG")

	AppConfig.REDISHost = viper.GetString("REDIS_HOST")
	AppConfig.REDISPassword = viper.GetString("REDIS_PASS")
	AppConfig.REDISDbno = viper.GetInt("REDIS_DBNO")
	AppConfig.REDISExpired = viper.GetInt("REDIS_EXPIRED")

	// check
	if AppConfig.Port == 0 || AppConfig.Host == "" || AppConfig.REDISHost == "" || AppConfig.REDISPassword == "" || AppConfig.REDISExpired == 0 {
		return errors.New("required variabel environment is empty")
	}

	log.Println("[INIT] configuration loaded")
	return nil
}
