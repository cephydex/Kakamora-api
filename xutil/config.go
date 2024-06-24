package xutil

import "github.com/spf13/viper"

type Config struct {
	AppPort      string `mapstructure:"APP_PORT"`
	SlackToken   string `mapstructure:"SLACK_TOKEN"`
	SlackChannel string `mapstructure:"SLACK_CHANNEL"`
	SendGridUrl   string `mapstructure:"SENDGRID_URL"`
	SendGridApiKey string `mapstructure:"SENDGRID_API_KEY"`

	FromEmail   string `mapstructure:"FROM_EMAIL"`
	FromName string `mapstructure:"FROM_NAME"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
