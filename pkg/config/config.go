package config

import "github.com/ilyakaznacheev/cleanenv"

type ConfigDatabase struct {
    Port     string `env:"POSTGRES_PORT" env-default:"5432"`
    Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
    Name     string `env:"POSTGRES_DB" env-default:"postgres"`
    User     string `env:"POSTGRES_USER" env-default:"user"`
    Password string `env:"POSTGRES_PASSWORD"`
}

type ConfigHttp struct{
    Port     string `env:"PORT" env-default:"5432"`
    Host     string `env:"HOST" env-default:"localhost"`
}

var cnfDB ConfigDatabase
var cnfHTTP ConfigHttp

func MustLoad() *ConfigDatabase {
  err := cleanenv.ReadConfig(".env", &cnfDB);
  if err != nil {
    panic("Error while get config")
  }
  err = cleanenv.ReadConfig(".env", &cnfHTTP);
  if err != nil {
    panic("Error while get config")
  }
  return &cnfDB
}
