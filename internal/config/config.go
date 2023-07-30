package config

import "github.com/ilyakaznacheev/cleanenv"

type ConfigDatabase struct {
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
	Name     string `env:"POSTGRES_DB" env-default:"postgres"`
	User     string `env:"POSTGRES_USER" env-default:"user"`
	Password string `env:"POSTGRES_PASSWORD"`
}

type ConfigHttp struct {
	Port string `env:"PORT" env-default:"8000"`
	Host string `env:"HOST" env-default:"localhost"`
}

type ConfigGrpc struct {
	Port string `env:"GRPC_PORT" env-default:"3333"`
	Host string `env:"GRPC_HOST" env-default:"localhost"`
}

type MainConfig struct {
	HttpConf     *ConfigHttp
	DataBaseConf *ConfigDatabase
	GrpcConf     *ConfigGrpc
}

var cnfDB ConfigDatabase
var cnfHTTP ConfigHttp
var cnfGRPC ConfigGrpc

func MustLoad() *MainConfig {
	err := cleanenv.ReadConfig(".env", &cnfDB)
	if err != nil {
		panic("Error while get config")
	}
	err = cleanenv.ReadConfig(".env", &cnfHTTP)
	if err != nil {
		panic("Error while get config")
	}
	err = cleanenv.ReadConfig(".env", &cnfGRPC)
	if err != nil {
		panic("Error while get config")
	}
	return &MainConfig{&cnfHTTP, &cnfDB, &cnfGRPC}
}
