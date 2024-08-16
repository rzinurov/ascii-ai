package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ImageGenerator ImageGeneratorConfig
	ImageStore     ImageStoreConfig
}

type ImageGeneratorConfig struct {
	Token  string
	Prompt string
}

type ImageStoreConfig struct {
	Dir string
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
		ImageGenerator: ImageGeneratorConfig{
			Token:  viper.GetString("image_generator.token"),
			Prompt: viper.GetString("image_generator.prompt"),
		},
		ImageStore: ImageStoreConfig{
			Dir: viper.GetString("image_store.dir"),
		},
	}

	return config
}
