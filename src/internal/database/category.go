package database

import (
	"database/sql"

	"github.com/gaspartv/api.ecommerce/src/internal/entity"
)

type CategoryDB struct {
	db *sql.DB
}

func NewCategoryDB(db *sql.DB) *CategoryDB {
	return &CategoryDB{db: db}
}

func (c *CategoryDB) Create(category *entity.CategoryCreate) {}
