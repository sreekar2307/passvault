package repository

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"passVault/dtos"
	"passVault/interfaces"
	"passVault/models"
)

type UserRepositoryImpl struct{}

func NewUserRepository() interfaces.UserRepository {
	return UserRepositoryImpl{}
}

func (u UserRepositoryImpl) GetUser(ctx context.Context, db *gorm.DB, filter dtos.GetUserFilter, user *models.User) error {
	return db.WithContext(ctx).Where(filter).Preload("UserSalt").First(user).Error
}

func (u UserRepositoryImpl) CreateUser(ctx context.Context, db *gorm.DB, user *models.User) error {
	return db.WithContext(ctx).Clauses(clause.Returning{}).Create(user).Error
}
