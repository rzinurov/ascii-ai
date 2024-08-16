package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	OpenAi      OpenAiConfig
	ImageFolder string
}

type OpenAiConfig struct {
	Token  string
	Prompt string
}

func Init() Config {
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal("Missing config.yaml file. Please create one from the provided example.")
		} else {
			log.Fatal("Unable to read config file", err)
		}
	}

	config := Config{
		OpenAi: OpenAiConfig{
			Token:  viper.GetString("openai.token"),
			Prompt: viper.GetString("openai.prompt"),
		},
		ImageFolder: viper.GetString("image_folder"),
	}

	return config
}
