package repository

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"passVault/dtos"
	"passVault/interfaces"
	"passVault/models"
)

type WebAuthnSessionRepositoryImpl struct{}

func NewWebAuthnSessionRepository() interfaces.WebAuthnSessionRepository {
	return WebAuthnSessionRepositoryImpl{}
}

func (w WebAuthnSessionRepositoryImpl) GetWebAuthnSession(ctx context.Context, db *gorm.DB, filter dtos.GetWebAuthnSessionFilter, session *models.WebauthNSession) error {
	return db.WithContext(ctx).Preload("User").Where(filter).First(session).Error
}

func (w WebAuthnSessionRepositoryImpl) CreateWebAuthnSession(ctx context.Context, db *gorm.DB, session *models.WebauthNSession) error {
	return db.WithContext(ctx).Clauses(clause.Returning{}).Create(session).Error
}
