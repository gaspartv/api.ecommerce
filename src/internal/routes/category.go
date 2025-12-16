package routes

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gaspartv/api.ecommerce/src/config"
	"github.com/gaspartv/api.ecommerce/src/internal/handler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CategoryRoutes(router *gin.Engine, db *gorm.DB, r2 *s3.Client, env *config.Env) {
	categoryHandler := handler.NewCategoryHandler(db, r2, env)
	categoryGroup := router.Group("categories")
	{
		categoryGroup.POST("create", categoryHandler.Create)
		categoryGroup.GET("list", categoryHandler.List)
		categoryGroup.GET("find", categoryHandler.GetByID)
		categoryGroup.PATCH("edit", categoryHandler.Edit)
		categoryGroup.DELETE("delete", categoryHandler.Delete)
		categoryGroup.PATCH("disable", categoryHandler.Disable)
		categoryGroup.PATCH("change-image", categoryHandler.ChangeImage)
	}
}
