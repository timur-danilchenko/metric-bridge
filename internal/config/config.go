package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Kafka struct {
		Brokers []string `mapstructure:"brokers"`
		Topic   string   `mapstructure:"topic"`
	} `mapstructure:"kafka"`

	Postgres struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Dbname   string `mapstructure:"dbname"`
	}

	Prometheus struct {
		Port int `mapstructure:"port"`
	}
}

func Load(path string) (config Config, err error) {
	// Set up viper to read the config.yaml file
	viper.SetConfigName("config") // Config file name without extension
	viper.SetConfigType("yaml")   // Config file type

	// Look for the config file in the configs directory
	viper.AddConfigPath("./configs")

	// Read the config file
	err = viper.ReadInConfig()
	if err != nil {
		return
		// log.Fatalf("Error reading config file: %s", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
		// log.Fatalf("Unable to decode into struct, %v", err)
	}
	return
}
