package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	OpenAi struct {
		Token string
	}
	ImageFolder string
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

	config := Config{}
	config.OpenAi.Token = viper.GetString("openai.token")
	config.ImageFolder = viper.GetString("image_folder")

	return config
}
