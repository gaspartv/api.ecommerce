package entity

import (
	"time"

	"github.com/gaspartv/api.ecommerce/src/config"
	"github.com/nrednav/cuid2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID              string         `gorm:"type:varchar(32);primaryKey" json:"id"`
	CreatedAt       time.Time      `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt       *time.Time     `gorm:"type:timestamptz" json:"updated_at,omitempty"`
	DisabledAt      *time.Time     `gorm:"type:timestamptz" json:"disabled_at,omitempty"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name            string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"name"`
	Image           string         `gorm:"type:varchar(255)" json:"image"`
	Email           string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	EmailVerifiedAt *time.Time     `gorm:"type:timestamptz" json:"email_verified_at,omitempty"`
	PasswordHash    string         `gorm:"type:varchar(255);not null" json:"-"`
}

type UserCreate struct {
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	PasswordHash string `json:"password" binding:"required"`
}

func NewUser(user UserCreate, env *config.Env) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), 10)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:           cuid2.Generate(),
		Name:         user.Name,
		Image:        env.IMAGE_CATEGORY_DEFAULT_URL,
		Email:        user.Email,
		PasswordHash: string(hash),
	}, nil
}
