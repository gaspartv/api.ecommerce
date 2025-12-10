package handler

import (
	"github.com/gaspartv/api.ecommerce/src/internal/entity"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryHandle struct {
	db *gorm.DB
}

func NewCategoryHandler(db *gorm.DB) *CategoryHandle {
	return &CategoryHandle{db: db}
}

func (h *CategoryHandle) Create(ctx *gin.Context) {
	var categoryCreate entity.CategoryCreate

	if err := ctx.BindJSON(&categoryCreate); err != nil {
		ctx.JSON(422, gin.H{"error": err.Error()})
		return
	}

	category := entity.NewCategory(categoryCreate)

	// categoryExists := h.db.Where("name = ?", category.Name).First(&entity.Category{})
	// if categoryExists.Error == nil {
	// 	ctx.JSON(409, gin.H{"error": "Category with this name already exists"})
	// 	return
	// }

	result := h.db.Create(category)
	if result.Error != nil {
		ctx.JSON(400, gin.H{"error": result.Error.Error()})
		return
	}

	ctx.JSON(201, gin.H{"data": category})
}

func (h *CategoryHandle) List(ctx *gin.Context) {
	categories := make([]entity.Category, 0)

	if err := h.db.Where("deleted_at IS NULL").Find(&categories).Error; err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"data": categories})
}
