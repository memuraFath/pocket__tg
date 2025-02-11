package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken     string
	PocketConsumerKey string
	AuthServerURL     string
	TelegramBotUrl    string `mapstructure:"bot_url"`
	DBPath            string `mapstructure:"db_file"`
	Messages          Messages
}

type Messages struct {
	Errors
	Responses
}

type Errors struct {
	InvalidUrl    string `mapstructure:"invalid_url"`
	NotAuthorized string `mapstructure:"not_authorized"`
	UnableToSave  string `mapstructure:"unable_to_save"`
	Default       string `mapstructure:"default"`
}

type Responses struct {
	Start             string `mapstructure:"start"`
	LinkSaved         string `mapstructure:"link_saved"`
	AlreadyAuthorized string `mapstructure:"already_authorozed"`
	UnknownCommand    string `mapstructure:"unknown_command"`
}

func Init() (*Config, error) {
	if err := setUpViper(); err != nil {
		return nil, err
	}

	var cfg Config

	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := parseEnv(&cfg); err != nil {
		return nil, err
	}

	log.Println(cfg)
	return &cfg, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.Unmarshal(cfg); err != nil {
		log.WithFields(log.Fields{
			"handler": "config.Unmarshal",
			"problem": "can not unmarshal config",
		}).Error(err)
		return err
	}
	//log.Println("1")
	if err := viper.UnmarshalKey("messages.responses", &cfg.Messages.Responses); err != nil {
		log.WithFields(log.Fields{
			"handler": "config.Unmarshal",
			"problem": "can not unmarshalKey cfg.messages.responses",
		}).Error(err)
		return err
	}

	log.Println("2")
	if err := viper.UnmarshalKey("messages.errors", &cfg.Messages.Errors); err != nil {
		return err
	}

	return nil
}

func setUpViper() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")
	return viper.ReadInConfig()
}

func parseEnv(cfg *Config) error {
	if err := viper.BindEnv("TOKEN"); err != nil {
		return err
	}

	cfg.TelegramToken = viper.GetString("TOKEN")

	if err := viper.BindEnv("CONSUMER_KEY"); err != nil {
		return err
	}

	cfg.PocketConsumerKey = viper.GetString("CONSUMER_KEY")

	if err := viper.BindEnv("AUTH_SERVER_URL"); err != nil {
		return err
	}
	cfg.AuthServerURL = viper.GetString("AUTH_SERVER_URL")

	return nil
}
