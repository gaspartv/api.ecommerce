package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gaspartv/api.ecommerce/src/config"
	"github.com/gaspartv/api.ecommerce/src/internal/entity"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductHandle struct {
	db  *gorm.DB
	env *config.Env
}

func NewProductHandler(db *gorm.DB, env *config.Env) *ProductHandle {
	return &ProductHandle{
		db:  db,
		env: env,
	}
}

func (h *ProductHandle) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(ctx.Query("limit"))
	if limit < 1 {
		limit = 20
	}

	orderBy := ctx.Query("order_by")
	if orderBy == "" {
		orderBy = "updated_at"
	}

	orderDir := ctx.Query("order_dir")
	if orderDir != "asc" && orderDir != "desc" {
		orderDir = "desc"
	}

	offset := (page - 1) * limit

	query := h.db.Model(&entity.Product{}).Where("deleted_at IS NULL")

	var products []entity.Product
	var total int64

	if err := query.Count(&total).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := query.
		Limit(limit).
		Offset(offset).
		Order(orderBy + " " + orderDir).
		Find(&products).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  products,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

func (h *ProductHandle) Create(ctx *gin.Context) {
	var productCreate entity.ProductCreate

	if err := ctx.BindJSON(&productCreate); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	product := entity.NewProduct(productCreate, h.env)

	if err := h.db.Create(&product).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": product})
}

func (h *ProductHandle) Edit(ctx *gin.Context) {
	idParam := ctx.Query("id")
	if idParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	var body entity.ProductEdit
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var current entity.Product
	if err := h.db.Where("id = ? AND deleted_at IS NULL", idParam).First(&current).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}

	if body.Name != nil {
		updates["name"] = *body.Name
	}

	if body.Description != nil {
		updates["description"] = *body.Description
	}

	if body.Price != nil {
		updates["price"] = *body.Price
	}

	if body.StockQuantity != nil {
		updates["stock_quantity"] = *body.StockQuantity
	}

	if body.CategoryID != nil {
		updates["category_id"] = *body.CategoryID
	}

	if body.Sku != nil {
		updates["sku"] = *body.Sku
	}

	if body.Weight != nil {
		updates["weight"] = *body.Weight
	}

	if body.Dimensions != nil {
		updates["dimensions"] = *body.Dimensions
	}

	if body.IsFeatured != nil {
		updates["is_featured"] = *body.IsFeatured
	}

	if len(updates) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	if err := h.db.Model(&entity.Product{}).Where("id = ? AND deleted_at IS NULL", idParam).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}
