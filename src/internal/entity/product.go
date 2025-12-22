package entity

import (
	"time"

	"github.com/gaspartv/api.ecommerce/src/config"
	"github.com/nrednav/cuid2"
	"gorm.io/gorm"
)

type Product struct {
	ID            string         `gorm:"type:varchar(32);primaryKey" json:"id"`
	CreatedAt     time.Time      `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt     *time.Time     `gorm:"type:timestamptz" json:"updated_at,omitempty"`
	DisabledAt    *time.Time     `gorm:"type:timestamptz" json:"disabled_at,omitempty"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name          string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"name"`
	Description   string         `gorm:"type:varchar(510)" json:"description"`
	Image         string         `gorm:"type:varchar(255)" json:"image"`
	Price         float64        `gorm:"type:numeric(10,2);not null" json:"price"`
	StockQuantity int            `gorm:"type:int;not null" json:"stock_quantity"`
	CategoryID    string         `gorm:"type:varchar(32);not null;index" json:"category_id"`
	Sku           string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"sku"`
	Weight        float64        `gorm:"type:numeric(10,3)" json:"weight"`
	Dimensions    string         `gorm:"type:varchar(100)" json:"dimensions"`
	IsFeatured    bool           `gorm:"type:boolean;not null;default:false" json:"is_featured"`
}

type ProductWithCategory struct {
	Product
	CategoryName string `json:"category_name"`
}

type ProductCreate struct {
	Name          string  `json:"name" binding:"required"`
	Description   string  `json:"description" binding:"required"`
	Price         float64 `json:"price" binding:"required"`
	StockQuantity int     `json:"stock_quantity" binding:"required"`
	CategoryID    string  `json:"category_id" binding:"required"`
	Sku           string  `json:"sku" binding:"required"`
	Weight        float64 `json:"weight,omitempty"`
	Dimensions    string  `json:"dimensions,omitempty"`
	IsFeatured    bool    `json:"is_featured,omitempty"`
}

type ProductEdit struct {
	Name          *string  `json:"name,omitempty"`
	Description   *string  `json:"description,omitempty"`
	Price         *float64 `json:"price,omitempty"`
	StockQuantity *int     `json:"stock_quantity,omitempty"`
	CategoryID    *string  `json:"category_id,omitempty"`
	Sku           *string  `json:"sku,omitempty"`
	Weight        *float64 `json:"weight,omitempty"`
	Dimensions    *string  `json:"dimensions,omitempty"`
	IsFeatured    *bool    `json:"is_featured,omitempty"`
}

func NewProduct(create ProductCreate, env *config.Env) *Product {
	return &Product{
		ID:            cuid2.Generate(),
		Name:          create.Name,
		Description:   create.Description,
		Image:         env.IMAGE_CATEGORY_DEFAULT_URL,
		Price:         create.Price,
		StockQuantity: create.StockQuantity,
		CategoryID:    create.CategoryID,
		Sku:           create.Sku,
		Weight:        create.Weight,
		Dimensions:    create.Dimensions,
		IsFeatured:    create.IsFeatured,
	}
}
