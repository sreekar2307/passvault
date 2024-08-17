package services

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"passVault/dtos"
	"passVault/interfaces"
	"passVault/models"
	"passVault/resources"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	db                           *gorm.DB
	encryptionService            interfaces.EncryptService
	hashService                  interfaces.HashService
	userRepository               interfaces.UserRepository
	userSaltRepository           interfaces.UserSaltRepository
	captchaService               interfaces.CaptchaService
	webAuthnCredentialRepository interfaces.WebAuthnCredentialRepository
	webAuthnSessionRepository    interfaces.WebAuthnSessionRepository
}

func NewUserService(
	db *gorm.DB,
	encryptionService interfaces.EncryptService,
	hashService interfaces.HashService,
	userRepository interfaces.UserRepository,
	userSaltRepository interfaces.UserSaltRepository,
	captchaService interfaces.CaptchaService,
	webAuthnCredentialRepository interfaces.WebAuthnCredentialRepository,
	webAuthnSessionRepository interfaces.WebAuthnSessionRepository,
) interfaces.UserService {
	return UserServiceImpl{
		db:                           db,
		encryptionService:            encryptionService,
		userRepository:               userRepository,
		userSaltRepository:           userSaltRepository,
		hashService:                  hashService,
		captchaService:               captchaService,
		webAuthnCredentialRepository: webAuthnCredentialRepository,
		webAuthnSessionRepository:    webAuthnSessionRepository,
	}
}

func (u UserServiceImpl) CreateUser(ctx context.Context, params dtos.CreateUserParams) (string, error) {
	var (
		user     models.User
		userSalt models.UserSalt
	)
	if ok, err := u.captchaService.VerifyToken(ctx, params.Token); err != nil {
		return "", err
	} else if !ok {
		return "", errors.New("invalid captcha")
	}

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

func (u UserServiceImpl) BeginWebAuthnRegister(ctx context.Context, params dtos.RegisterWebAuthnUserParams) (dtos.BeginRegisterResult, error) {
	var (
		user     models.User
		userSalt models.UserSalt
	)
	if err := u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := u.userRepository.GetUser(ctx, tx, dtos.GetUserFilter{Email: params.Email}, &user); err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			salt, err := u.encryptionService.GenerateSalt(ctx)
			if err != nil {
				return err
			}
			userSalt = models.UserSalt{
				Salt: salt,
			}
			user = models.User{
				Name:  params.Name,
				Email: params.Email,
			}
			if err := u.userSaltRepository.CreateUserSalt(ctx, tx, &userSalt); err != nil {
				return err
			}
			user.UserSaltID = userSalt.ID
			if err := u.userRepository.CreateUser(ctx, tx, &user); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return dtos.BeginRegisterResult{}, err
	}
	if user.IsRegistered {
		return dtos.BeginRegisterResult{}, errors.New("user already registered")
	}
	webAuthN := resources.WebAuthn()
	creation, session, err := webAuthN.BeginRegistration(user)
	if err != nil {
		return dtos.BeginRegisterResult{}, err
	}
	webAuthNSession := models.WebauthNSession{
		UserID: user.ID,
		SessionData: models.SessionData{
			Challenge:        session.Challenge,
			UserID:           base64.RawURLEncoding.EncodeToString(session.UserID),
			UserVerification: session.UserVerification,
			Extensions:       session.Extensions,
		},
	}

	if !session.Expires.IsZero() {
		webAuthNSession.SessionData.Expires = &session.Expires
	}

	for _, cred := range session.AllowedCredentialIDs {
		webAuthNSession.SessionData.AllowedCredentialIDs = append(webAuthNSession.SessionData.AllowedCredentialIDs,
			base64.RawURLEncoding.EncodeToString(cred))
	}

	err = u.webAuthnSessionRepository.CreateWebAuthnSession(ctx, u.db, &webAuthNSession)
	if err != nil {
		return dtos.BeginRegisterResult{}, err
	}
	resp := dtos.BeginRegisterResult{
		CredOptions: creation,
		SessionID:   webAuthNSession.SessionID.String(),
	}
	return resp, nil
}

func (u UserServiceImpl) FinishWebAuthnRegister(ctx context.Context, sessionIDStr string, reqBody io.Reader) (string, error) {
	var (
		logger  = resources.Logger(ctx)
		session models.WebauthNSession
	)
	webAuthN := resources.WebAuthn()
	pcc, err := protocol.ParseCredentialCreationResponseBody(reqBody)
	if err != nil {
		var typedErr *protocol.Error
		errors.As(err, &typedErr)
		logger.Error("error parsing credential creation response body",
			"type", typedErr.Type, "details", typedErr.Details, "error", err)
		return "", err
	}
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return "", err
	}
	if err := u.webAuthnSessionRepository.GetWebAuthnSession(ctx, u.db,
		dtos.GetWebAuthnSessionFilter{SessionID: sessionID},
		&session); err != nil {
		return "", err
	}
	sessionData := webauthn.SessionData{
		Challenge:        session.SessionData.Challenge,
		UserVerification: session.SessionData.UserVerification,
		Extensions:       session.SessionData.Extensions,
	}
	userID, err := base64.RawURLEncoding.DecodeString(session.SessionData.UserID)
	if err != nil {
		return "", err
	}
	sessionData.UserID = userID
	for _, cred := range session.SessionData.AllowedCredentialIDs {
		credAsBytes, err := base64.RawURLEncoding.DecodeString(cred)
		if err != nil {
			return "", err
		}
		sessionData.AllowedCredentialIDs = append(sessionData.AllowedCredentialIDs, credAsBytes)
	}
	if session.SessionData.Expires != nil {
		sessionData.Expires = *session.SessionData.Expires
	}

	cred, err := webAuthN.CreateCredential(session.User, sessionData, pcc)
	if err != nil {
		return "", err
	}

	if err := u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := u.webAuthnCredentialRepository.CreateWebAuthnCredential(ctx, tx, &models.WebauthNCredential{
			UserID: session.UserID,
			Credential: models.Credential{
				ID:              base64.RawURLEncoding.EncodeToString(cred.ID),
				PublicKey:       base64.RawURLEncoding.EncodeToString(cred.PublicKey),
				AttestationType: cred.AttestationType,
				Transport:       cred.Transport,
				Flags:           cred.Flags,
				Authenticator:   cred.Authenticator,
			},
		}); err != nil {
			return err
		}
		return u.userRepository.UpdateUser(ctx, tx, dtos.GetUserFilter{ID: session.UserID},
			&models.User{IsRegistered: true})
	}); err != nil {
		return "", err
	}

	token, err := u.newToken(ctx, session.User)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u UserServiceImpl) Login(ctx context.Context, params dtos.LoginParams) (string, error) {
	var user models.User

	if ok, err := u.captchaService.VerifyToken(ctx, params.Token); err != nil {
		return "", err
	} else if !ok {
		return "", errors.New("invalid captcha")
	}

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

func (u UserServiceImpl) BeginWebAuthnLogin(ctx context.Context, params dtos.BeginLoginParams) (dtos.BeginLoginResult, error) {
	var (
		user        models.User
		credentials []models.WebauthNCredential
	)

	if err := u.userRepository.GetUser(ctx, u.db, dtos.GetUserFilter{Email: params.Email}, &user); err != nil {
		return dtos.BeginLoginResult{}, err
	}
	if err := u.webAuthnCredentialRepository.GetWebAuthnCredentials(ctx, u.db, dtos.GetWebAuthnCredentialFilter{UserID: user.ID}, &credentials); err != nil {
		return dtos.BeginLoginResult{}, err
	}
	for _, cred := range credentials {
		webauthnCred := webauthn.Credential{
			AttestationType: cred.Credential.AttestationType,
			Transport:       cred.Credential.Transport,
			Flags:           cred.Credential.Flags,
			Authenticator:   cred.Credential.Authenticator,
		}
		webauthnCred.ID, _ = base64.RawURLEncoding.DecodeString(cred.Credential.ID)
		webauthnCred.PublicKey, _ = base64.RawURLEncoding.DecodeString(cred.Credential.PublicKey)
		user.Credentials = append(user.Credentials, webauthnCred)
	}
	webAuthn := resources.WebAuthn()
	options, session, err := webAuthn.BeginLogin(user)
	if err != nil {
		return dtos.BeginLoginResult{}, err
	}
	webAuthNSession := models.WebauthNSession{
		UserID: user.ID,
		SessionData: models.SessionData{
			Challenge:        session.Challenge,
			UserID:           base64.RawURLEncoding.EncodeToString(session.UserID),
			UserVerification: session.UserVerification,
			Extensions:       session.Extensions,
		},
	}

	if !session.Expires.IsZero() {
		webAuthNSession.SessionData.Expires = &session.Expires
	}

	for _, cred := range session.AllowedCredentialIDs {
		webAuthNSession.SessionData.AllowedCredentialIDs = append(webAuthNSession.SessionData.AllowedCredentialIDs,
			base64.RawURLEncoding.EncodeToString(cred))
	}

	err = u.webAuthnSessionRepository.CreateWebAuthnSession(ctx, u.db, &webAuthNSession)
	if err != nil {
		return dtos.BeginLoginResult{}, err
	}
	return dtos.BeginLoginResult{
		CredAssertion: options,
		SessionID:     webAuthNSession.SessionID.String(),
	}, nil
}

func (u UserServiceImpl) FinishWebAuthnLogin(ctx context.Context, sessionIDStr string, reqBody io.Reader) (string, error) {
	var (
		session     models.WebauthNSession
		credentials []models.WebauthNCredential
	)

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return "", err
	}
	if err := u.webAuthnSessionRepository.GetWebAuthnSession(ctx, u.db,
		dtos.GetWebAuthnSessionFilter{SessionID: sessionID},
		&session); err != nil {
		return "", err
	}

	if err := u.webAuthnCredentialRepository.GetWebAuthnCredentials(ctx,
		u.db, dtos.GetWebAuthnCredentialFilter{UserID: session.User.ID}, &credentials); err != nil {
		return "", err
	}
	for _, cred := range credentials {
		webauthnCred := webauthn.Credential{
			AttestationType: cred.Credential.AttestationType,
			Transport:       cred.Credential.Transport,
			Flags:           cred.Credential.Flags,
			Authenticator:   cred.Credential.Authenticator,
		}
		webauthnCred.ID, _ = base64.RawURLEncoding.DecodeString(cred.Credential.ID)
		webauthnCred.PublicKey, _ = base64.RawURLEncoding.DecodeString(cred.Credential.PublicKey)
		session.User.Credentials = append(session.User.Credentials, webauthnCred)
	}

	webAuthn := resources.WebAuthn()
	par, err := protocol.ParseCredentialRequestResponseBody(reqBody)
	if err != nil {
		return "", err
	}
	sessionData := webauthn.SessionData{
		Challenge:        session.SessionData.Challenge,
		UserVerification: session.SessionData.UserVerification,
		Extensions:       session.SessionData.Extensions,
	}
	userID, err := base64.RawURLEncoding.DecodeString(session.SessionData.UserID)
	if err != nil {
		return "", err
	}
	sessionData.UserID = userID
	for _, cred := range session.SessionData.AllowedCredentialIDs {
		credAsBytes, err := base64.RawURLEncoding.DecodeString(cred)
		if err != nil {
			return "", err
		}
		sessionData.AllowedCredentialIDs = append(sessionData.AllowedCredentialIDs, credAsBytes)
	}
	if session.SessionData.Expires != nil {
		sessionData.Expires = *session.SessionData.Expires
	}
	credential, err := webAuthn.ValidateLogin(session.User, sessionData, par)
	if err != nil {
		return "", err
	}

	for _, cred := range credentials {
		if cred.Credential.ID == base64.RawURLEncoding.EncodeToString(credential.ID) {
			if err := u.webAuthnCredentialRepository.UpdateWebAuthnCredential(ctx, u.db,
				dtos.GetWebAuthnCredentialFilter{CredID: cred.CredID}, &models.WebauthNCredential{
					Credential: models.Credential{
						ID:              base64.RawURLEncoding.EncodeToString(credential.ID),
						PublicKey:       base64.RawURLEncoding.EncodeToString(credential.PublicKey),
						AttestationType: credential.AttestationType,
						Transport:       credential.Transport,
						Flags:           credential.Flags,
						Authenticator:   credential.Authenticator,
					},
				}); err != nil {
				return "", err
			}
			break
		}
	}

	token, err := u.newToken(ctx, session.User)
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
