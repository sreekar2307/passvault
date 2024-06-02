package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"passVault/interfaces"
	"passVault/models"
)

type PasswordVersionRepositoryImpl struct{}

func NewPasswordVersionRepository() interfaces.PasswordVersionRepository {
	return PasswordVersionRepositoryImpl{}
}

func (p PasswordVersionRepositoryImpl) InsertNewVersion(
	ctx context.Context,
	db *gorm.DB,
	password models.Password,
) error {
	var latestVersion models.PasswordVersion
	if err := db.Where(models.PasswordVersion{PasswordID: password.ID}).Last(&latestVersion).Error; err != nil &&
		!errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	passwordVersion := models.PasswordVersion{
		PasswordID:     password.ID,
		Username:       password.Username,
		Notes:          password.Notes,
		Password:       password.Password,
		HashedPassword: password.HashedPassword,
		Website:        password.Website,
		Email:          password.Email,
		Version:        latestVersion.Version + 1,
	}
	return db.Clauses(clause.Returning{}).Create(&passwordVersion).Error

}
