package services

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"passVault/dtos"
	"passVault/interfaces"
	"passVault/models"
	"passVault/resources"
	"time"
)

type UserServiceImpl struct {
	db                 *gorm.DB
	encryptionService  interfaces.EncryptService
	hashService        interfaces.HashService
	userRepository     interfaces.UserRepository
	userSaltRepository interfaces.UserSaltRepository
}

func NewUserService(
	db *gorm.DB,
	encryptionService interfaces.EncryptService,
	hashService interfaces.HashService,
	userRepository interfaces.UserRepository,
	userSaltRepository interfaces.UserSaltRepository,

) interfaces.UserService {
	return UserServiceImpl{
		db:                 db,
		encryptionService:  encryptionService,
		userRepository:     userRepository,
		userSaltRepository: userSaltRepository,
		hashService:        hashService,
	}
}

func (u UserServiceImpl) CreateUser(ctx context.Context, params dtos.CreateUserParams) (string, error) {
	var (
		user     models.User
		userSalt models.UserSalt
	)
	if err := u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		salt, err := u.encryptionService.GenerateSalt(ctx)
		if err != nil {
			return err
		}
		userSalt = models.UserSalt{
			Salt: salt,
		}
		passwordHash, err := u.hashService.Hash(ctx, params.Password)
		if err != nil {
			return err
		}
		user = models.User{
			Name:     params.Name,
			Email:    params.Email,
			Password: passwordHash,
		}

		if err := u.userSaltRepository.CreateUserSalt(ctx, tx, &userSalt); err != nil {
			return err
		}
		user.UserSaltID = userSalt.ID
		if err := u.userRepository.CreateUser(ctx, tx, &user); err != nil {
			return err
		}
		return nil

	}); err != nil {
		return "", err
	}
	token, err := u.newToken(ctx, user)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u UserServiceImpl) Login(ctx context.Context, params dtos.LoginParams) (string, error) {
	var user models.User

	if err := u.userRepository.GetUser(ctx, u.db, dtos.GetUserFilter{Email: params.Email}, &user); err != nil {
		return "", err
	}
	if ok, err := u.hashService.CompareHash(ctx, user.Password, params.Password); err != nil {
		return "", errors.New("incorrect password")
	} else if !ok {
		return "", nil
	}
	token, err := u.newToken(ctx, user)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u UserServiceImpl) newToken(_ context.Context, user models.User) (string, error) {
	var (
		now        = time.Now()
		exp        = now.Add(time.Hour * 24).Unix()
		config     = resources.Config()
		encryption dtos.Encryption
	)
	claims := jwt.MapClaims{
		"exp":    exp,
		"iat":    now.Unix(),
		"nbf":    now.Unix(),
		"userID": user.ID,
		"iss":    "passVault",
	}
	if err := config.UnmarshalKey(dtos.ConfigKeys.Encryption, &encryption); err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	shaKey, err := hex.DecodeString(encryption.Auth)
	if err != nil {
		return "", err
	}
	token, err := jwtToken.SignedString(shaKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u UserServiceImpl) ValidateToken(ctx context.Context, token string, user *models.User) error {
	var (
		config     = resources.Config()
		encryption dtos.Encryption
		claims     jwt.MapClaims
	)
	if err := config.UnmarshalKey(dtos.ConfigKeys.Encryption, &encryption); err != nil {
		return err
	}
	_, err := jwt.NewParser(
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithTimeFunc(func() time.Time { return time.Now() }),
		jwt.WithLeeway(time.Minute),
		jwt.WithIssuer("passVault"),
	).ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return hex.DecodeString(encryption.Auth)
	})
	if err != nil {
		return err
	}
	userIDAsFloat, ok := claims["userID"].(float64)
	if !ok {
		return errors.New("invalid token")
	}
	userID := uint(userIDAsFloat)
	return u.userRepository.GetUser(ctx, u.db, dtos.GetUserFilter{ID: userID}, user)
}
