package model

type Configuration struct {
	Network  string
	Address  string
	Port     int
	Database DatabaseConfiguration
	Password Password
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
}
