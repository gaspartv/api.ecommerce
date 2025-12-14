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

type CategoryHandle struct {
	db  *gorm.DB
	r2  *s3.Client
	env *config.Env
}

func NewCategoryHandler(db *gorm.DB, r2 *s3.Client, env *config.Env) *CategoryHandle {
	return &CategoryHandle{
		db:  db,
		r2:  r2,
		env: env,
	}
}

func (h *CategoryHandle) Create(ctx *gin.Context) {
	var categoryCreate entity.CategoryCreate

	if err := ctx.BindJSON(&categoryCreate); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	category := entity.NewCategory(categoryCreate, h.env)

	err := h.db.Where("name = ?", category.Name).First(&entity.Category{}).Error
	if err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Category already exists"})
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := h.db.Create(category)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": category})
}

func (h *CategoryHandle) List(ctx *gin.Context) {
	search := ctx.Query("search")

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

	query := h.db.Model(&entity.Category{}).Where("deleted_at IS NULL")

	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", like, like)
	}

	var categories []entity.Category
	var total int64

	if err := query.Count(&total).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := query.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", orderBy, orderDir)).
		Find(&categories).Error; err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  categories,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *CategoryHandle) GetByID(ctx *gin.Context) {
	idParam := ctx.Query("id")
	if idParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	category := &entity.Category{}

	result := h.db.Where("id = ? AND deleted_at IS NULL", idParam).First(category)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": category})
}

func (h *CategoryHandle) Edit(ctx *gin.Context) {
	idParam := ctx.Query("id")
	if idParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	var body entity.CategoryEdit
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var current entity.Category
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
		newName := *body.Name
		if newName != current.Name {
			err := h.db.Where("name = ? AND id <> ? AND deleted_at IS NULL", newName, idParam).First(&entity.Category{}).Error
			if err == nil {
				ctx.JSON(http.StatusConflict, gin.H{"error": "Category name already exists"})
				return
			}

			if !errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			updates["name"] = newName
		}
	}

	if body.Description != nil && *body.Description != current.Description {
		updates["description"] = *body.Description
	}

	if len(updates) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	if err := h.db.Model(&entity.Category{}).Where("id = ? AND deleted_at IS NULL", idParam).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Category updated successfully"})
}

func (h *CategoryHandle) Delete(ctx *gin.Context) {
	idParam := ctx.Query("id")
	if idParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	result := h.db.Where("id = ? AND deleted_at IS NULL", idParam).Delete(&entity.Category{})
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

func (h *CategoryHandle) Disable(ctx *gin.Context) {
	idParam := ctx.Query("id")
	if idParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	var cat entity.Category

	if err := h.db.Where("id = ? AND deleted_at IS NULL", idParam).First(&cat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var newValue interface{}
	if cat.DisabledAt == nil {
		newValue = gorm.Expr("NOW()")
	} else {
		newValue = nil
	}

	if err := h.db.Model(&entity.Category{}).
		Where("id = ?", idParam).
		Update("disabled_at", newValue).Error; err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	msg := "Category enabled successfully"
	if cat.DisabledAt == nil {
		msg = "Category disabled successfully"
	}

	ctx.JSON(http.StatusOK, gin.H{"message": msg})
}

func (h *CategoryHandle) ChangeImage(ctx *gin.Context) {
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

	result := h.db.Model(&entity.Category{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("image", url)

	if result.Error != nil {
		ctx.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(404, gin.H{"error": "Category not found"})
		return
	}

	ctx.JSON(200, gin.H{"message": "Image updated successfully"})
}
