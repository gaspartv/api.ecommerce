package handler

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gaspartv/api.ecommerce/src/config"
	"github.com/gaspartv/api.ecommerce/src/internal/entity"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandle struct {
	db  *gorm.DB
	r2  *s3.Client
	env *config.Env
}

func NewUserHandler(db *gorm.DB, r2 *s3.Client, env *config.Env) *UserHandle {
	return &UserHandle{
		db:  db,
		r2:  r2,
		env: env,
	}
}

func (h *UserHandle) Create(ctx *gin.Context) {
	var userCreate entity.UserCreate

	if err := ctx.BindJSON(&userCreate); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user, err := entity.NewUser(userCreate, h.env)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": user})
}
