package model

type Configuration struct {
	Network  string
	Address  string
	Port     int
	Database DatabaseConfiguration
}

type DatabaseConfiguration struct {
	Host     string
	Port     int
	UserName string
	Password string
	DbName   string
}
