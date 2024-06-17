package repository

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"passVault/dtos"
	"passVault/interfaces"
	"passVault/models"
)

type WebAuthnCredentialRepositoryImpl struct{}

func NewWebAuthnCredentialRepository() interfaces.WebAuthnCredentialRepository {
	return WebAuthnCredentialRepositoryImpl{}
}

func (w WebAuthnCredentialRepositoryImpl) GetWebAuthnCredentials(ctx context.Context, db *gorm.DB, filter dtos.GetWebAuthnCredentialFilter, credentials *[]models.WebauthNCredential) error {
	return db.WithContext(ctx).Where(filter).Find(credentials).Error
}

func (w WebAuthnCredentialRepositoryImpl) UpdateWebAuthnCredential(ctx context.Context, db *gorm.DB, filter dtos.GetWebAuthnCredentialFilter, credential *models.WebauthNCredential) error {
	return db.WithContext(ctx).Clauses(clause.Returning{}).Where(filter).Updates(credential).Error
}

func (w WebAuthnCredentialRepositoryImpl) CreateWebAuthnCredential(ctx context.Context, db *gorm.DB, credential *models.WebauthNCredential) error {
	return db.WithContext(ctx).Clauses(clause.Returning{}).Create(credential).Error
}
