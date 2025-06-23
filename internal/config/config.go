package config

import (
	"github.com/spf13/viper"
)

type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
	Sslmode  string
}

type PrometheusConfig struct {
	Port int `mapstructure:"port"`
}

type Config struct {
	Kafka      KafkaConfig      `mapstructure:"kafka"`
	Postgres   PostgresConfig   `mapstructure:"postgres"`
	Prometheus PrometheusConfig `mapstructure:"prometheus"`
}

func Load(path string) (cfg Config, err error) {
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

	viper.AutomaticEnv() // подгрузит из .env

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return
		// log.Fatalf("Unable to decode into struct, %v", err)
	}

	cfg.Postgres.User = viper.GetString("DB_USER")
	cfg.Postgres.Password = viper.GetString("DB_PASSWORD")
	cfg.Postgres.Host = viper.GetString("DB_HOST")
	cfg.Postgres.Port = viper.GetInt("DB_PORT")
	cfg.Postgres.Dbname = viper.GetString("DB_NAME")
	cfg.Postgres.Sslmode = viper.GetString("DB_SSLMODE")
	return
}
