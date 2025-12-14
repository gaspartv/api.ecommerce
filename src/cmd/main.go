package main

import (
	"fmt"
	"log"

	"github.com/gaspartv/api.ecommerce/src/config"
	"github.com/gaspartv/api.ecommerce/src/external/storage"
	"github.com/gaspartv/api.ecommerce/src/internal/entity"
	"github.com/gaspartv/api.ecommerce/src/internal/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

func main() {
	env, err := config.LoadEnv()
	if err != nil {
		log.Fatal("Erro ao carregar variáveis de ambiente:", err)
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

	r2, err := storage.NewR2Client(env)
	if err != nil {
		log.Fatal("Erro ao conectar no R2:", err)
	}

	router := gin.Default()

	// Configurar CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * 60 * 60,
	}))

	routes.CategoryRoutes(router, db, r2)

	router.Run(":" + env.Port)
}
