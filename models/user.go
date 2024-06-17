package models

import (
	"encoding/base64"
	"github.com/go-webauthn/webauthn/webauthn"
	"strconv"
)

type User struct {
	Base
	Password     string                `json:"password,omitempty" gorm:"column:password;type:text;not null"`
	IsRegistered bool                  `json:"is_registered,omitempty" gorm:"column:is_registered;type:boolean;not null;default:false"`
	UserSalt     UserSalt              `json:"-" gorm:"foreignKey:UserSaltID"`
	UserSaltID   uint                  `json:"-" gorm:"column:user_salt_id;type:bigint;not null"`
	Name         string                `json:"name,omitempty" gorm:"column:name;type:text;not null;default:''"`
	Email        string                `json:"email,omitempty" gorm:"column:email;type:text;not null;unique"`
	Credentials  []webauthn.Credential `json:"-" gorm:"-"`
}

func (u User) WebAuthnID() []byte {
	res, _ := base64.RawURLEncoding.DecodeString(strconv.Itoa(int(u.ID)))
	return res
}

func (u User) WebAuthnName() string {
	return u.Name
}

func (u User) WebAuthnDisplayName() string {
	return u.Name
}

func (u User) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}

func (u User) WebAuthnIcon() string {
	return ""
}
