package routes

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gaspartv/api.ecommerce/src/config"
	"github.com/gaspartv/api.ecommerce/src/internal/handler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ProductRoutes(router *gin.Engine, db *gorm.DB, r2 *s3.Client, env *config.Env) {
	productHandler := handler.NewProductHandler(db, r2, env)
	productGroup := router.Group("products")
	{
		productGroup.POST("create", productHandler.Create)
		productGroup.GET("list", productHandler.List)
		productGroup.PATCH("edit", productHandler.Edit)
		productGroup.DELETE("delete", productHandler.Delete)
		productGroup.PATCH("disable", productHandler.Disable)
		productGroup.PATCH("change-image", productHandler.ChangeImage)
	}
}
