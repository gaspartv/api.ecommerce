package entity

import (
	"time"

	"github.com/nrednav/cuid2"
)

type Category struct {
	ID          string     `gorm:"type:varchar(32);primaryKey"`
	CreatedAt   time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt   *time.Time `gorm:"type:timestamptz"`
	DisabledAt  *time.Time `gorm:"type:timestamptz"`
	DeletedAt   *time.Time `gorm:"type:timestamptz"`
	Name        string     `gorm:"type:varchar(255);not null;uniqueIndex"`
	Description string     `gorm:"type:varchar(510)"`
	Image       string     `gorm:"type:varchar(255)"`
}

type CategoryCreate struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Image       string `json:"image" binding:"required"`
}

func NewCategory(create CategoryCreate) *Category {
	return &Category{
		ID:          cuid2.Generate(),
		Name:        create.Name,
		Description: create.Description,
		Image:       create.Image,
	}
}
