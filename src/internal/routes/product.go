package routes

import (
	"github.com/gaspartv/api.ecommerce/src/config"
	"github.com/gaspartv/api.ecommerce/src/internal/handler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ProductRoutes(router *gin.Engine, db *gorm.DB, env *config.Env) {
	productHandler := handler.NewProductHandler(db, env)
	productGroup := router.Group("products")
	{
		productGroup.POST("create", productHandler.Create)
		productGroup.GET("list", productHandler.List)
		productGroup.PATCH("edit", productHandler.Edit)
		// productGroup.DELETE("delete", productHandler.Delete)
		// productGroup.PATCH("disable", productHandler.Disable)
		// productGroup.PATCH("change-image", productHandler.ChangeImage)
	}
}
