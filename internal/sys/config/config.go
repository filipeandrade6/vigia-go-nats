package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Auth struct {
	Directory string `mapstructure:"directory"`
	ActiveKID string `mapstructure:"activekid"`
}

type Database struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Name         string `mapstructure:"name"`
	MaxIDLEConns int    `mapstructure:"maxidleconns"`
	MaxOpenConns int    `mapstructure:"maxopenconns"`
	// SSLMode      string `mapstructure:"sslmode"`
	DisableTLS bool `mapstructure:"disabletls"`
}

type Service struct {
	Host string `mapstructure:"host"`
	Conn string `mapstructure:"conn"`
	Port int    `mapstructure:"port"`
}

type Configuration struct {
	Auth     Auth     `mapstructure:"auth"`
	Database Database `mapstructure:"database"`
	Service  Service  `mapstructure:"service"`
}

func ParseConfig(build string) (Configuration, error) {
	viper.SetDefault("auth.directory", "deployments/keys")
	viper.SetDefault("auth.activekid", "bcc18baa-7830-4cfc-8f96-8a26ede5d81f")
	viper.SetDefault("database.host", "dev-postgres:5432")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "secret")
	viper.SetDefault("database.name", "vigia")
	viper.SetDefault("database.maxidleconns", "0")
	viper.SetDefault("database.maxopenconns", "0")
	// viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.disabletls", "true")
	viper.SetDefault("service.host", "gerencia")
	viper.SetDefault("service.conn", "tcp")
	viper.SetDefault("service.port", "12346")

	viper.BindEnv("auth.directory", "VIGIA_AUTH_DIR")
	viper.BindEnv("auth.activekid", "VIGIA_AUTH_ACTIVEKID")
	viper.BindEnv("database.host", "VIGIA_DB_HOST")
	viper.BindEnv("database.user", "VIGIA_DB_USER")
	viper.BindEnv("database.password", "VIGIA_DB_PASSWORD")
	viper.BindEnv("database.name", "VIGIA_DB_NAME")
	viper.BindEnv("database.maxidleconns", "VIGIA_DB_MAXIDLECONNS")
	viper.BindEnv("database.maxopenconns", "VIGIA_DB_MAXOPENCONNS")
	// viper.BindEnv("database.sslmode", "VIGIA_DB_SSLMODE")
	viper.BindEnv("database.disabletls", "VIGIA_DB_DISABLETLS")
	viper.BindEnv("service.host", "VIGIA_GER_HOST")
	viper.BindEnv("service.conn", "VIGIA_GER_SERVER_CONN")
	viper.BindEnv("service.port", "VIGIA_GER_SERVER_PORT")

	viper.AutomaticEnv()
	fmt.Println(viper.Get("service.host"))

	cfg := Configuration{}

	if err := viper.Unmarshal(&cfg); err != nil {
		return Configuration{}, fmt.Errorf("unmsarshalling config: %w", err)
	}

	return cfg, nil
}

// TODO melhorar isso aqui
func PrettyPrintConfig() string {
	return "entrou no prettyPrintConfig"
	// return fmt.Sprintf(
	// 	"VIGIA_GER_HOST: %s, VIGIA_GER_SERVER_CONN: %s, VIGIA_GER_SERVER_PORT: %d, VIGIA_GRA_HOST: %s, VIGIA_GRA_CONN: %s, VIGIA_GRA_PORT: %d, VIGIA_DB_HOST: %s, VIGIA_DB_USER: %s, VIGIA_DB_PASSWORD: %s, VIGIA_DB_NAME: %s, VIGIA_DB_MAXIDLECONNS: %d, VIGIA_DB_MAXOPENCONNS: %d, VIGIA_DB_SSLMODE: %s, VIGIA_MET_WEB_DEBUGHOST: %s, VIGIA_MET_EXPVAR_HOST: %s, VIGIA_MET_EXPVAR_ROUTE: %s, VIGIA_MET_EXPVAR_READTIMEOUT: %d, VIGIA_MET_EXPVAR_WRITETIMEOUT: %d, VIGIA_MET_EXPVAR_IDLETIMEOUT: %d, VIGIA_MET_EXPVAR_SHUTDOWNTIMEOUT: %d, VIGIA_MET_COLLECT_FROM: %s, VIGIA_MET_PUBLISH_TO: %s, VIGIA_MET_PUBLISH_INTERVAL: %d",

	// 	viper.GetString("VIGIA_GER_HOST"),
	// 	viper.GetString("VIGIA_GER_SERVER_CONN"),
	// 	viper.GetInt("VIGIA_GER_SERVER_PORT"),

	// 	viper.GetString("VIGIA_GRA_HOST"),
	// 	viper.GetString("VIGIA_GRA_CONN"),
	// 	viper.GetInt("VIGIA_GRA_PORT"),

	// 	viper.GetString("VIGIA_DB_HOST"),
	// 	viper.GetString("VIGIA_DB_USER"),
	// 	viper.GetString("VIGIA_DB_PASSWORD"),
	// 	viper.GetString("VIGIA_DB_NAME"),
	// 	viper.GetInt("VIGIA_DB_MAXIDLECONNS"),
	// 	viper.GetInt("VIGIA_DB_MAXOPENCONNS"),
	// 	viper.GetString("VIGIA_DB_SSLMODE"),

	// 	viper.GetString("VIGIA_MET_WEB_DEBUGHOST"),
	// 	viper.GetString("VIGIA_MET_EXPVAR_HOST"),
	// 	viper.GetString("VIGIA_MET_EXPVAR_ROUTE"),
	// 	viper.GetInt("VIGIA_MET_EXPVAR_READTIMEOUT"),
	// 	viper.GetInt("VIGIA_MET_EXPVAR_WRITETIMEOUT"),
	// 	viper.GetInt("VIGIA_MET_EXPVAR_IDLETIMEOUT"),
	// 	viper.GetInt("VIGIA_MET_EXPVAR_SHUTDOWNTIMEOUT"),
	// 	viper.GetString("VIGIA_MET_COLLECT_FROM"),
	// 	viper.GetString("VIGIA_MET_PUBLISH_TO"),
	// 	viper.GetInt("VIGIA_MET_PUBLISH_INTERVAL"),
	// )
}
