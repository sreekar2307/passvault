package interfaces

import (
	"context"
	"gorm.io/gorm"
	"passVault/dtos"
	"passVault/models"
)

type PasswordRepository interface {
	GetPassword(context.Context, *gorm.DB, dtos.GetPasswordFilter, *models.Password) error
	DeletePassword(context.Context, *gorm.DB, dtos.GetPasswordFilter) error
	GetPasswords(context.Context, *gorm.DB, dtos.GetPasswordFilter, *[]models.Password) error
	CreatePassword(context.Context, *gorm.DB, *models.Password) error
	CreatePasswords(context.Context, *gorm.DB, []*models.Password) error
	UpdatePassword(context.Context, *gorm.DB, dtos.GetPasswordFilter, *models.Password) error
}

type PasswordGenerationRepository interface {
	GenerationHistory(context.Context, *gorm.DB, *models.PasswordGenerationHistory) error
}

type PasswordVersionRepository interface {
	InsertNewVersion(context.Context, *gorm.DB, models.Password) error
}

type UserRepository interface {
	GetUser(context.Context, *gorm.DB, dtos.GetUserFilter, *models.User) error
	CreateUser(context.Context, *gorm.DB, *models.User) error
}

type UserSaltRepository interface {
	CreateUserSalt(context.Context, *gorm.DB, *models.UserSalt) error
}
