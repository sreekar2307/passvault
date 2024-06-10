package services

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io"
	"passVault/dtos"
	"passVault/interfaces"
	"passVault/models"
	"passVault/resources"
	"passVault/utils"
)

type PasswordServiceImpl struct {
	db                           *gorm.DB
	passwordRepository           interfaces.PasswordRepository
	passwordGenerationRepository interfaces.PasswordGenerationRepository
	passwordVersionRepository    interfaces.PasswordVersionRepository
	encryptionService            interfaces.EncryptService
	hashService                  interfaces.HashService
}

func NewPasswordService(
	db *gorm.DB,
	passwordRepository interfaces.PasswordRepository,
	passwordGenerationRepository interfaces.PasswordGenerationRepository,
	passwordVersionRepository interfaces.PasswordVersionRepository,
	encryptionService interfaces.EncryptService,
	hashService interfaces.HashService,

) interfaces.PasswordService {
	return PasswordServiceImpl{
		db:                           db,
		passwordRepository:           passwordRepository,
		passwordGenerationRepository: passwordGenerationRepository,
		passwordVersionRepository:    passwordVersionRepository,
		encryptionService:            encryptionService,
		hashService:                  hashService,
	}
}

func (p PasswordServiceImpl) GetPasswords(ctx context.Context, user models.User, params dtos.GetPasswordsParams) ([]models.Password, error) {
	var passwords []models.Password
	if err := p.passwordRepository.GetPasswords(ctx, p.db, dtos.GetPasswordFilter{
		UserID:      user.ID,
		NameLike:    params.Query,
		Limit:       params.Limit,
		Offset:      params.Offset,
		WebsiteLike: params.Query,
		Email:       params.Query,
	}, &passwords); err != nil {
		return nil, err
	}

	for i := range passwords {
		var err error
		passwords[i].Password, err = p.encryptionService.Decrypt(ctx, user.UserSalt, passwords[i].Password)
		if err != nil {
			return nil, err
		}
	}

	return passwords, nil
}

func (p PasswordServiceImpl) GetPassword(ctx context.Context, user models.User, id uint) (models.Password, error) {
	var password models.Password
	if err := p.passwordRepository.GetPassword(ctx, p.db, dtos.GetPasswordFilter{
		UserID: user.ID,
		ID:     id,
	}, &password); err != nil {
		return models.Password{}, err
	}

	var err error

	password.Password, err = p.encryptionService.Decrypt(ctx, user.UserSalt, password.Password)

	return password, err
}

func (p PasswordServiceImpl) DeletePassword(ctx context.Context, user models.User, id uint) (bool, error) {
	if err := p.passwordRepository.DeletePassword(ctx, p.db, dtos.GetPasswordFilter{
		UserID: user.ID,
		ID:     id,
	}); err != nil {
		return false, err
	}
	return true, nil
}

func (p PasswordServiceImpl) GeneratePassword(ctx context.Context, user models.User, params dtos.GeneratePasswordParams) (string, error) {
	var sampleSpace []rune
	if !params.ExcludeAlphabets {
		sampleSpace = append(sampleSpace, utils.Alpha...)
	}
	if !params.ExcludeDigits {
		sampleSpace = append(sampleSpace, utils.Numeric...)
	}
	if !params.ExcludeSymbols {
		sampleSpace = append(sampleSpace, utils.Symbols...)
	}
	if len(sampleSpace) == 0 {
		return "", errors.New("no sample space provided")
	}
	password := utils.RandFromSampleSpace(params.Size, sampleSpace)

	if err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		passwordEncrypted, err := p.encryptionService.Encrypt(ctx, user.UserSalt, password)
		if err != nil {
			return fmt.Errorf("encrypt password: %w", err)
		}

		if err := p.passwordGenerationRepository.GenerationHistory(ctx, tx, &models.PasswordGenerationHistory{
			Password: passwordEncrypted,
			UserID:   user.ID,
		}); err != nil {
			return fmt.Errorf("store password generation history: %w", err)
		}
		return nil
	}); err != nil {
		return "", err
	}
	return password, nil
}

func (p PasswordServiceImpl) StorePassword(ctx context.Context, user models.User, params dtos.StorePasswordParams) (bool, error) {

	if err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		password := models.Password{
			Name:     utils.SqlNullStringIfValidFromString(params.Name),
			Password: params.Password,
			UserID:   user.ID,
			Notes:    utils.SqlNullStringIfValidFromString(params.Notes),
			Username: utils.SqlNullStringIfValidFromString(params.Username),
			Email:    utils.SqlNullStringIfValidFromString(params.Email),
			Website:  params.Website,
		}
		err := p.createPassword(ctx, tx, password, user)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return false, err
	}
	return true, nil
}

func (p PasswordServiceImpl) createPassword(ctx context.Context,
	tx *gorm.DB,
	password models.Password,
	user models.User,
) error {

	err := p.encryptPassword(ctx, user, &password)
	if err != nil {
		return err
	}
	if err := p.passwordRepository.CreatePassword(ctx, tx, &password); err != nil {
		return fmt.Errorf("store password: %w", err)
	}
	if err := p.passwordVersionRepository.InsertNewVersion(ctx, tx, password); err != nil {
		return fmt.Errorf("insert new version: %w", err)
	}
	return nil
}

func (p PasswordServiceImpl) encryptPassword(ctx context.Context, user models.User, password *models.Password) error {
	encryptedPlainPassword, err := p.encryptionService.Encrypt(ctx, user.UserSalt, password.Password)
	if err != nil {
		return fmt.Errorf("encrypt password: %w", err)
	}

	hashedPassword, err := p.hashService.Hash(ctx, password.Password)
	if err != nil {
		return err
	}
	password.Password = encryptedPlainPassword
	password.HashedPassword = hashedPassword
	return nil
}

func (p PasswordServiceImpl) UpdatePassword(ctx context.Context, user models.User, params dtos.UpdatePasswordParams) (bool, error) {
	var (
		password        models.Password
		updatesPassword models.Password
	)
	if err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := p.passwordRepository.GetPassword(ctx, tx, dtos.GetPasswordFilter{
			ID:     params.ID,
			UserID: user.ID,
		}, &password); err != nil {
			return fmt.Errorf("get password: %w", err)
		}
		updatesPassword.Name = utils.SqlNullStringIfValidFromString(params.Name)
		updatesPassword.Website = params.Website
		updatesPassword.Username = utils.SqlNullStringIfValidFromString(params.Username)
		updatesPassword.Email = utils.SqlNullStringIfValidFromString(params.Email)
		updatesPassword.Notes = utils.SqlNullStringIfValidFromString(params.Notes)

		if params.Password != "" {
			updatesPassword.Password = params.Password
			if err := p.encryptPassword(ctx, user, &updatesPassword); err != nil {
				return err
			}
		}

		if err := p.passwordRepository.UpdatePassword(ctx, tx,
			dtos.GetPasswordFilter{ID: password.ID},
			&updatesPassword); err != nil {
			return fmt.Errorf("update password: %w", err)
		}

		if err := p.passwordVersionRepository.InsertNewVersion(ctx, tx, updatesPassword); err != nil {
			return fmt.Errorf("insert new version: %w", err)
		}

		return nil
	}); err != nil {
		return false, err
	}
	return true, nil
}

func (p PasswordServiceImpl) ImportPasswords(ctx context.Context, user models.User, csvFile io.ReadCloser) (int, error) {
	defer func() {
		_ = csvFile.Close()
	}()
	csvReader := csv.NewReader(csvFile)
	var (
		passwordsInserted int
		logger            = resources.Logger(ctx)
	)
	_, err := csvReader.Read()
	if err != nil {
		return 0, err
	}
	for record, err := csvReader.Read(); err != io.EOF; record, err = csvReader.Read() {
		var (
			username = record[0]
			name     = record[3]
			password = record[4]
			website  = record[6]
		)
		if _, err := p.StorePassword(ctx, user, dtos.StorePasswordParams{
			Username: username,
			Email:    username,
			Password: password,
			Website:  website,
			Name:     name,
		}); err != nil {
			logger.Error("failed to store password", "error", err.Error(), "record", record)
		} else {
			passwordsInserted++
		}
	}
	return passwordsInserted, nil
}
