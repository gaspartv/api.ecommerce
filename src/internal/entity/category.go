package entity

import (
	"time"

	"github.com/nrednav/cuid2"
	"gorm.io/gorm"
)

type Category struct {
	ID          string         `gorm:"type:varchar(32);primaryKey" json:"id"`
	CreatedAt   time.Time      `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt   *time.Time     `gorm:"type:timestamptz" json:"updated_at,omitempty"`
	DisabledAt  *time.Time     `gorm:"type:timestamptz" json:"disabled_at,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name        string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"name"`
	Description string         `gorm:"type:varchar(510)" json:"description"`
	Image       string         `gorm:"type:varchar(255)" json:"image"`
}

type CategoryCreate struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type CategoryEdit struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type CategoryChangeImage struct {
	Image string `json:"image" binding:"required"`
}

func NewCategory(create CategoryCreate) *Category {
	return &Category{
		ID:          cuid2.Generate(),
		Name:        create.Name,
		Description: create.Description,
	}
}
