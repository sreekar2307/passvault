package repository

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"passVault/interfaces"
	"passVault/models"
)

type PasswordGenerationRepositoryImpl struct{}

func NewPasswordGenerationRepository() interfaces.PasswordGenerationRepository {
	return PasswordGenerationRepositoryImpl{}
}

func (p PasswordGenerationRepositoryImpl) GenerationHistory(
	ctx context.Context,
	db *gorm.DB,
	passwordGenerationHistory *models.PasswordGenerationHistory,
) error {
	return db.WithContext(ctx).Clauses(clause.Returning{}).Create(passwordGenerationHistory).Error
}
