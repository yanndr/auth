package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

// AppSettings represent the settings for the application.
type AppSettings struct {
	Network   string
	Address   string
	GRPCPort  int
	TLSConfig TLS
	Database  Database
	Password  Password
	Token     Token
}

// TLS settings
type TLS struct {
	UseTLS        bool
	CertFile      string
	KeyFile       string
	CAFile        string
	ServerAddress string
}

// Database settings
type Database struct {
	Type     string
	Path     string
	Host     string
	Port     int
	UserName string
	Password string
	DbName   string
	SslMode  string
	RootCert string
	SslKey   string
	SslCert  string
}

// Password settings
type Password struct {
	MinLength    int
	MinNumeric   int
	MinUpperCase int
	MinLowerCase int
	MinSpecial   int
}

// Token settings
type Token struct {
	SigningMethod string
	SignedKey     string
	Audience      string
	Issuer        string
	ExpDuration   int
}

// LoadConfiguration parse a file (configName) Json or Yaml in the path configPath and returns an AppSettings configuration struct
func LoadConfiguration(configName, configPath string) (*AppSettings, error) {
	configuration := &AppSettings{}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
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
