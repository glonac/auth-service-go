package config

import "github.com/ilyakaznacheev/cleanenv"

type ConfigDatabase struct {
	Port     string `env:"POSTGRES_PORT_AUTH" env-default:"5432"`
	Host     string `env:"POSTGRES_HOST_AUTH" env-default:"localhost"`
	Name     string `env:"POSTGRES_DB_AUTH" env-default:"postgres"`
	User     string `env:"POSTGRES_USER_AUTH" env-default:"user"`
	Password string `env:"POSTGRES_PASSWORD_AUTH"`
}

type ConfigHttp struct {
	Port string `env:"PORT_AUTH_SERVICE" env-default:"8001"`
	Host string `env:"HOST_AUTH_SERVICE" env-default:"localhost"`
}

type ConfigGrpc struct {
	Port string `env:"GRPC_PORT" env-default:"3333"`
	Host string `env:"GRPC_HOST" env-default:"localhost"`
}

type ConfigQueue struct {
	Port     string `env:"RABBITMQ_PORT" env-default:"5672"`
	Host     string `env:"RABBITMQ_HOST" env-default:"localhost"`
	User     string `env:"RABBITMQ_DEFAULT_USER" env-default:"guest"`
	Password string `env:"RABBITMQ_DEFAULT_PASS" env-default:"guest"`
	Vhost    string `env:"RABBITMQ_DEFAULT_VHOST" env-default:"/"`
}

type MainConfig struct {
	HttpConf     *ConfigHttp
	DataBaseConf *ConfigDatabase
	GrpcConf     *ConfigGrpc
	QueueConf    *ConfigQueue
}

var cnfDB ConfigDatabase
var cnfHTTP ConfigHttp
var cnfGRPC ConfigGrpc
var cnfQueue ConfigQueue

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
	err = cleanenv.ReadConfig(".env", &cnfQueue)
	if err != nil {
		panic("Error while get config")
	}
	return &MainConfig{&cnfHTTP, &cnfDB, &cnfGRPC, &cnfQueue}
}
