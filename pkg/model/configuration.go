package model

type Configuration struct {
	Network   string
	Address   string
	GRPCPort  int
	TLSConfig TLSConfig
	Database  DatabaseConfiguration
	Password  Password
	Token     Token
}

type TLSConfig struct {
	UseTLS        bool
	CertFile      string
	KeyFile       string
	CAFile        string
	ServerAddress string
}

type DatabaseConfiguration struct {
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
	SignedKey   string
	ExpDuration int
}
