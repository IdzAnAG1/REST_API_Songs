package config

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	DataBaseURL    string
	APIPort        string
	ExternalApiURL string
}

func LoadConfig() (*Config, error) {
	if os.Getenv("DATABASE_URL") == "" ||
		os.Getenv("API_PORT") == "" ||
		os.Getenv("EXTERNAL_API_URL") == "" {
		logrus.Warn("Environment variables not found uploading .env file")
		if err := godotenv.Load(); err != nil {
			logrus.Fatal("File .env is not exist. Please verify its existence. " +
				"If you have difficulties, please follow this link  " +
				"https://github.com/IdzAnAG1/REST_API_Songs ")
		}
	}
	config := &Config{
		DataBaseURL:    os.Getenv("DATABASE_URL"),
		APIPort:        os.Getenv("API_PORT"),
		ExternalApiURL: os.Getenv("EXTERNAL_API_URL"),
	}
	if config.DataBaseURL == "" || config.ExternalApiURL == "" || config.APIPort == " " {
		logrus.Fatal("One or more environment variables are not set." +
			" Check the environment variables or the .env file")
	}
	logrus.Info("Configuration is loaded")
	return config, nil
}
