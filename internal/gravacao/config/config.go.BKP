package config

import (
	"fmt"
	// "os"
	// "path/filepath"

	"github.com/spf13/viper"
)

// type Auth struct {
// 	Directory string `mapstructure:"directory"`
// 	ActiveKID string `mapstructure:"activekid"`
// }

type Gravacao struct {
	// Conn          string `mapstructure:"conn"`
	Port int `mapstructure:"port"`
	// Armazenamento string `mapstructure:"armazenamento"`
	// Housekeeper   int    `mapstructure:"housekeeper"`
}

type Configuration struct {
	Build string
	// Auth     Auth     `mapstructure:"auth"`
	Gravacao Gravacao `mapstructure:"gravacao"`
}

func ParseConfig(build string) (Configuration, error) {
	// userDir, err := os.UserHomeDir()
	// if err != nil {
	// 	return Configuration{}, fmt.Errorf("getting user directory: %w", err)
	// }

	// viper.SetDefault("auth.directory", "deployments/keys")
	// viper.SetDefault("auth.activekid", "bcc18baa-7830-4cfc-8f96-8a26ede5d81f")
	// viper.SetDefault("gravacao.conn", "tcp")
	viper.SetDefault("gravacao.port", "12346")
	// viper.SetDefault("gravacao.armazenamento", filepath.Join(userDir, "vigia"))
	// viper.SetDefault("gravacao.housekeeper", "168")

	// viper.BindEnv("auth.directory", "VIGIA_AUTH_DIR")
	// viper.BindEnv("auth.activekid", "VIGIA_AUTH_ACTIVEKID")
	// viper.BindEnv("gravacao.conn", "VIGIA_GRA_SERVER_CONN")
	viper.BindEnv("gravacao.port", "VIGIA_GRA_SERVER_PORT")
	// viper.BindEnv("gravacao.armazenamento", "VIGIA_GRA_SERVER_ARMAZENAMENTO")
	// viper.BindEnv("gravacao.housekeeper", "VIGIA_GRA_SERVER_HOUSEKEEPER")

	viper.AutomaticEnv()

	cfg := Configuration{Build: build}

	if err := viper.Unmarshal(&cfg); err != nil {
		return Configuration{}, fmt.Errorf("unmsarshalling config: %w", err)
	}

	return cfg, nil
}
