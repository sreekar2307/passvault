package repository

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"passVault/dtos"
	"passVault/interfaces"
	"passVault/models"
)

type PasswordRepositoryImpl struct{}

func NewPasswordRepository() interfaces.PasswordRepository {
	return PasswordRepositoryImpl{}
}

func (p PasswordRepositoryImpl) GetPassword(ctx context.Context, db *gorm.DB, filter dtos.GetPasswordFilter, password *models.Password) error {
	return db.WithContext(ctx).Where(filter).First(password).Error
}

func (p PasswordRepositoryImpl) DeletePassword(ctx context.Context, db *gorm.DB, filter dtos.GetPasswordFilter) error {
	return db.WithContext(ctx).Where(filter).Delete(&models.Password{}).Error
}

func (p PasswordRepositoryImpl) GetPasswords(ctx context.Context, db *gorm.DB, filter dtos.GetPasswordFilter, passwords *[]models.Password) error {
	query := db.WithContext(ctx).Where(models.Password{
		UserID: filter.UserID,
		Base:   models.Base{ID: filter.ID},
	})
	var (
		searchQuery   = db.WithContext(ctx)
		showAddSearch = false
	)
	if filter.NameLike != "" {
		searchQuery = searchQuery.Or("name ILIKE ?", "%"+filter.NameLike+"%")
		showAddSearch = true
	}
	if filter.WebsiteLike != "" {
		searchQuery = searchQuery.Or("website ILIKE ?", "%"+filter.WebsiteLike+"%")
		showAddSearch = true
	}
	if filter.Email != "" {
		searchQuery = searchQuery.Or("email ILIKE ?", "%"+filter.Email+"%")
		showAddSearch = true
	}
	if showAddSearch {
		query = query.Where(searchQuery)
	}
	if filter.Offset != 0 {
		query = query.Offset(filter.Offset)
	}
	if filter.Limit != 0 {
		query = query.Limit(filter.Limit)
	}
	return query.Order(clause.OrderByColumn{Column: clause.Column{Table: "passwords", Name: "id"}, Desc: true}).
		Find(passwords).Error
}

func (p PasswordRepositoryImpl) CreatePassword(ctx context.Context, db *gorm.DB, password *models.Password) error {
	return db.WithContext(ctx).Clauses(clause.Returning{}).Create(password).Error
}

func (p PasswordRepositoryImpl) CreatePasswords(ctx context.Context, db *gorm.DB, passwords []*models.Password) error {
	return db.WithContext(ctx).Clauses(clause.Returning{}).Create(passwords).Error
}

func (p PasswordRepositoryImpl) UpdatePassword(ctx context.Context, db *gorm.DB, filter dtos.GetPasswordFilter, password *models.Password) error {
	return db.WithContext(ctx).Clauses(clause.Returning{}).Where(filter).Updates(password).Error
}
