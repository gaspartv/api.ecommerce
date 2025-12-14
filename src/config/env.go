package config

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Env struct {
	Port                       string `validate:"required"`
	DatabaseHost               string `validate:"required"`
	DatabaseUser               string `validate:"required"`
	DatabasePass               string `validate:"required"`
	DatabaseName               string `validate:"required"`
	DatabasePort               string `validate:"required"`
	R2AccessKey                string `validate:"required"`
	R2SecretKey                string `validate:"required"`
	R2Endpoint                 string `validate:"required"`
	R2Bucket                   string `validate:"required"`
	R2PublicURL                string `validate:"required"`
	IMAGE_CATEGORY_DEFAULT_URL string `validate:"required"`
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
	env.R2AccessKey = os.Getenv("R2_ACCESS_KEY_ID")
	env.R2SecretKey = os.Getenv("R2_SECRET_ACCESS_KEY")
	env.R2Endpoint = os.Getenv("R2_ENDPOINT")
	env.R2Bucket = os.Getenv("R2_BUCKET")
	env.R2PublicURL = os.Getenv("R2_PUBLIC_URL")
	env.IMAGE_CATEGORY_DEFAULT_URL = os.Getenv("IMAGE_CATEGORY_DEFAULT_URL")

	validate := validator.New()
	if err := validate.Struct(env); err != nil {
		return nil, err
	}

	return &env, nil
}
