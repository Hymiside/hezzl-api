package models

type ConfigPostgres struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type ConfigClickhouse struct {
	Host string
	Port string
	Database string
}

type ConfigServer struct {
	Host string
	Port string
}

type ConfigRedis struct {
	Host string
	Port string
}

type ConfigNats struct {
	Host string
	Port string
}
