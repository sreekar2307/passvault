package repository

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"passVault/interfaces"
	"passVault/models"
)

type UserSaltRepositoryImpl struct{}

func NewUserSaltRepository() interfaces.UserSaltRepository {
	return UserSaltRepositoryImpl{}
}

func (u UserSaltRepositoryImpl) CreateUserSalt(ctx context.Context, db *gorm.DB, salt *models.UserSalt) error {
	return db.WithContext(ctx).Clauses(clause.Returning{}).Create(salt).Error
}
