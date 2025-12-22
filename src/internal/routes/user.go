package routes

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gaspartv/api.ecommerce/src/config"
	"github.com/gaspartv/api.ecommerce/src/internal/handler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(router *gin.Engine, db *gorm.DB, r2 *s3.Client, env *config.Env) {
	userHandler := handler.NewUserHandler(db, r2, env)
	userGroup := router.Group("users")
	{
		userGroup.POST("create", userHandler.Create)
		// userGroup.GET("list", userHandler.List)
		// userGroup.PATCH("edit", userHandler.Edit)
		// userGroup.DELETE("delete", userHandler.Delete)
		// userGroup.PATCH("disable", userHandler.Disable)
		// userGroup.PATCH("change-profile-image", userHandler.ChangeProfileImage)
	}
}
