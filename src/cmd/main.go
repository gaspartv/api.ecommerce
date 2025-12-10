package main

import (
	"fmt"
	"log"

	"github.com/gaspartv/api.ecommerce/src/config"
	"github.com/gaspartv/api.ecommerce/src/internal/entity"
	"github.com/gaspartv/api.ecommerce/src/internal/routes"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

func main() {
	env, err := config.LoadEnv()
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		env.DatabaseHost,
		env.DatabaseUser,
		env.DatabasePass,
		env.DatabaseName,
		env.DatabasePort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Erro ao conectar:", err)
	}
	db.AutoMigrate(&entity.Category{})

	router := gin.Default()

	routes.CategoryRoutes(router, db)

	router.Run(":" + env.Port)
}
