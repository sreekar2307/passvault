package interfaces

import (
	"context"
	"io"
	"passVault/dtos"
	"passVault/models"
)

type PasswordService interface {
	GeneratePassword(context.Context, models.User, dtos.GeneratePasswordParams) (string, error)
	StorePassword(context.Context, models.User, dtos.StorePasswordParams) (bool, error)
	UpdatePassword(context.Context, models.User, dtos.UpdatePasswordParams) (bool, error)
	DeletePassword(context.Context, models.User, uint) (bool, error)
	GetPasswords(context.Context, models.User, dtos.GetPasswordsParams) ([]models.Password, error)
	GetPassword(context.Context, models.User, uint) (models.Password, error)
	ImportPasswords(context.Context, models.User, io.ReadCloser) (int, error)
}

type EncryptService interface {
	Encrypt(context.Context, models.UserSalt, string) (string, error)
	EncryptBulk(context.Context, models.UserSalt, map[string]string) (map[string]string, error)
	GenerateSalt(context.Context) (string, error)
	Decrypt(context.Context, models.UserSalt, string) (string, error)
}

type HashService interface {
	Hash(context.Context, string) (string, error)
	CompareHash(context.Context, string, string) (bool, error)
}

type UserService interface {
	CreateUser(context.Context, dtos.CreateUserParams) (string, error)
	Login(context.Context, dtos.LoginParams) (string, error)
	ValidateToken(context.Context, string, *models.User) error
}

type BackupService interface {
	BackupDb(context.Context) error
}
