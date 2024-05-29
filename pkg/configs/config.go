package configs

import (
	"fmt"
	"os"

	configs "github.com/namnv2496/go-coffee-shop-demo/config"
	"gopkg.in/yaml.v2"
)

type ConfigFilePath string

type Config struct {
	GRPC     GRPC     `yaml:"grpc"`
	Database Database `yaml:"database"`
	Kafka    Kafka    `yaml:"kafka"`
	Redis    Redis    `yaml:"redis"`
	Cron     Cron     `yaml:"cron"`
	S3       S3       `yaml:"s3"`
}

func GetConfigFromYaml(filePath ConfigFilePath) (Config, error) {
	var (
		configBytes = configs.DefaultConfigBytes
		config      = Config{}
		err         error
	)

	if filePath != "" {
		configBytes, err = os.ReadFile(string(filePath))
		if err != nil {
			return Config{}, fmt.Errorf("failed to read YAML file: %w", err)
		}
	}

	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return config, nil
}
