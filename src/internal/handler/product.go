package handler

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gaspartv/api.ecommerce/src/config"
	"github.com/gaspartv/api.ecommerce/src/external/storage"
	"github.com/gaspartv/api.ecommerce/src/internal/entity"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductHandle struct {
	db  *gorm.DB
	r2  *s3.Client
	env *config.Env
}

func NewProductHandler(db *gorm.DB, r2 *s3.Client, env *config.Env) *ProductHandle {
	return &ProductHandle{
		db:  db,
		r2:  r2,
		env: env,
	}
}

func (h *ProductHandle) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(ctx.Query("limit"))
	fmt.Println(limit)
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

	query := h.db.Model(&entity.Product{}).Where("products.deleted_at IS NULL")

	search := ctx.Query("search")
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("products.name ILIKE ? OR products.description ILIKE ?", like, like)
	}

	status := ctx.Query("status")
	if status != "" {
		switch status {
		case "active":
			query = query.Where("products.disabled_at IS NULL")
		case "inactive":
			query = query.Where("products.disabled_at IS NOT NULL")
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
			return
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var products []entity.ProductWithCategory

	if err := query.
		Select("products.*, categories.name as category_name").
		Joins("LEFT JOIN categories ON products.category_id = categories.id").
		Limit(limit).
		Offset(offset).
		Order("products." + orderBy + " " + orderDir).
		Scan(&products).Error; err != nil {
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

func (h *ProductHandle) Delete(ctx *gin.Context) {
	idParam := ctx.Query("id")
	if idParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	result := h.db.Where("id = ? AND deleted_at IS NULL", idParam).Delete(&entity.Product{})
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func (h *ProductHandle) Disable(ctx *gin.Context) {
	idParam := ctx.Query("id")
	if idParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	var product entity.Product

	if err := h.db.Where("id = ? AND deleted_at IS NULL", idParam).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var newValue interface{}
	if product.DisabledAt == nil {
		newValue = gorm.Expr("NOW()")
	} else {
		newValue = nil
	}

	if err := h.db.Model(&entity.Product{}).
		Where("id = ?", idParam).
		Update("disabled_at", newValue).Error; err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	msg := "Product enabled successfully"
	status := "active"
	if product.DisabledAt == nil {
		msg = "Product disabled successfully"
		status = "inactive"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  status,
		"message": msg,
	})
}

func (h *ProductHandle) ChangeImage(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		ctx.JSON(400, gin.H{"error": "ID is required"})
		return
	}

	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid file"})
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	key := fmt.Sprintf("categories/%s%s", id, ext)

	url, err := storage.UploadToR2(h.r2, key, header.Header.Get("Content-Type"), file)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	result := h.db.Model(&entity.Product{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("image", url)

	if result.Error != nil {
		ctx.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(404, gin.H{"error": "Product not found"})
		return
	}

	ctx.JSON(200, gin.H{"message": "Image updated successfully"})
}
