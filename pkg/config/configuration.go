package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Application struct {
	Network   string
	Address   string
	GRPCPort  int
	TLSConfig TLS
	Database  Database
	Password  Password
	Token     Token
}

type TLS struct {
	UseTLS        bool
	CertFile      string
	KeyFile       string
	CAFile        string
	ServerAddress string
}

type Database struct {
	Host     string
	Port     int
	UserName string
	Password string
	DbName   string
}

type Password struct {
	MinLength    int
	MinNumeric   int
	MinUpperCase int
	MinLowerCase int
	MinSpecial   int
}

type Token struct {
	SigningMethod string
	SignedKey     string
	Audience      string
	Issuer        string
	ExpDuration   int
}

func LoadConfiguration(configName, configPath string) (*Application, error) {
	configuration := &Application{}
	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return configuration, fmt.Errorf("error reading config file: %w", err)
	}

	err := viper.Unmarshal(configuration)
	if err != nil {
		return configuration, fmt.Errorf("error reading config file: %w", err)
	}

	return configuration, nil
}
