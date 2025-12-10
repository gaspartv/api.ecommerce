package config

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Env struct {
	Port         string `validate:"required"`
	DatabaseHost string `validate:"required"`
	DatabaseUser string `validate:"required"`
	DatabasePass string `validate:"required"`
	DatabaseName string `validate:"required"`
	DatabasePort string `validate:"required"`
}

func LoadEnv() (*Env, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	var env Env
	env.Port = os.Getenv("PORT")
	env.DatabaseHost = os.Getenv("DATABASE_HOST")
	env.DatabaseUser = os.Getenv("DATABASE_USER")
	env.DatabasePass = os.Getenv("DATABASE_PASSWORD")
	env.DatabaseName = os.Getenv("DATABASE_NAME")
	env.DatabasePort = os.Getenv("DATABASE_PORT")

	validate := validator.New()
	if err := validate.Struct(env); err != nil {
		return nil, err
	}

	return &env, nil
}
