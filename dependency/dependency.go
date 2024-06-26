package dependency

import (
	"gorm.io/gorm"
	"net/http"
	"passVault/dtos"
	"passVault/interfaces"
	"passVault/repository"
	"passVault/resources"
	"passVault/services"
	"sync"
)

var (
	dependencies Dependency
	once         sync.Once
)

type Dependency struct {
	DB                     *gorm.DB
	Config                 interfaces.Config
	PasswordRepo           interfaces.PasswordRepository
	PasswordGenerationRepo interfaces.PasswordGenerationRepository
	PasswordVersionRepo    interfaces.PasswordVersionRepository
	UserSaltRepo           interfaces.UserSaltRepository
	UserRepo               interfaces.UserRepository
	WepAuthnCredentialRepo interfaces.WebAuthnCredentialRepository
	WepAuthnSessionRepo    interfaces.WebAuthnSessionRepository

	PasswordService interfaces.PasswordService
	UserService     interfaces.UserService
	EncryptService  interfaces.EncryptService
	HashService     interfaces.HashService
	BackupService   interfaces.BackupService
	CaptchaService  interfaces.CaptchaService
}

func init() {
	once.Do(func() {
		dependencies = newDependencies()
	})
}

func Dependencies() Dependency {
	return dependencies
}

func newDependencies() Dependency {
	var (
		config     = resources.Config()
		db         = resources.Database()
		encryption dtos.Encryption
	)

	dependencies = Dependency{
		PasswordRepo:           repository.NewPasswordRepository(),
		PasswordGenerationRepo: repository.NewPasswordGenerationRepository(),
		PasswordVersionRepo:    repository.NewPasswordVersionRepository(),
		WepAuthnCredentialRepo: repository.NewWebAuthnCredentialRepository(),
		WepAuthnSessionRepo:    repository.NewWebAuthnSessionRepository(),
		UserSaltRepo:           repository.NewUserSaltRepository(),
		UserRepo:               repository.NewUserRepository(),
		DB:                     db,
		Config:                 config,
		HashService:            services.NewHashService(),
	}

	if err := config.UnmarshalKey(dtos.ConfigKeys.Encryption, &encryption); err != nil {
		panic(err.Error())
	}
	dependencies.CaptchaService = services.NewCaptchaService(http.DefaultClient)
	dependencies.EncryptService = services.NewEncryptionService(db, encryption)
	dependencies.PasswordService = services.NewPasswordService(db, dependencies.PasswordRepo,
		dependencies.PasswordGenerationRepo, dependencies.PasswordVersionRepo,
		dependencies.EncryptService, dependencies.HashService)
	dependencies.UserService = services.NewUserService(db, dependencies.EncryptService, dependencies.HashService,
		dependencies.UserRepo, dependencies.UserSaltRepo, dependencies.CaptchaService,
		dependencies.WepAuthnCredentialRepo, dependencies.WepAuthnSessionRepo)

	dependencies.BackupService = services.NewBackupService(
		resources.NewS3(config.GetString(dtos.ConfigKeys.DatabaseBackup.Region)),
	)

	return dependencies
}
