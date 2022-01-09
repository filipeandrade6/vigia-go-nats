package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Database struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Name         string `mapstructure:"name"`
	MaxIDLEConns int    `mapstructure:"maxidleconns"`
	MaxOpenConns int    `mapstructure:"maxopenconns"`
	DisableTLS   bool   `mapstructure:"disabletls"`
}

type Configuration struct {
	Build    string
	Database Database `mapstructure:"database"`
}

func ParseConfig(build string) (Configuration, error) {
	viper.SetDefault("database.host", "dev-postgres:5432")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "secret")
	viper.SetDefault("database.name", "vigia")
	viper.SetDefault("database.maxidleconns", "0")
	viper.SetDefault("database.maxopenconns", "0")
	viper.SetDefault("database.disabletls", "true")

	viper.BindEnv("database.host", "VIGIA_DB_HOST")
	viper.BindEnv("database.user", "VIGIA_DB_USER")
	viper.BindEnv("database.password", "VIGIA_DB_PASSWORD")
	viper.BindEnv("database.name", "VIGIA_DB_NAME")
	viper.BindEnv("database.maxidleconns", "VIGIA_DB_MAXIDLECONNS")
	viper.BindEnv("database.maxopenconns", "VIGIA_DB_MAXOPENCONNS")
	viper.BindEnv("database.sslmode", "VIGIA_DB_DISABLETLS")

	viper.AutomaticEnv()

	cfg := Configuration{Build: build}

	if err := viper.Unmarshal(&cfg); err != nil {
		return Configuration{}, fmt.Errorf("unmsarshalling config: %w", err)
	}

	return cfg, nil
}
