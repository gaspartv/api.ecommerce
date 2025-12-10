package routes

import (
	"github.com/gaspartv/api.ecommerce/src/internal/handler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CategoryRoutes(router *gin.Engine, db *gorm.DB) {

	categoryHandler := handler.NewCategoryHandler(db)
	categoryGroup := router.Group("categories")
	{
		categoryGroup.POST("create", categoryHandler.Create)
		categoryGroup.GET("list", categoryHandler.List)
	}
}
