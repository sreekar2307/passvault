package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"time"
)

type WebauthNCredential struct {
	CreatedAt  time.Time  `json:"createdAt" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt  time.Time  `json:"updatedAt" gorm:"column:updated_at;type:timestamp;not null"`
	CredID     uuid.UUID  `gorm:"column:cred_id;default:uuid_generate_v4();type:uuid;not null;primaryKey" json:"credID"`
	UserID     uint       `json:"user_id "gorm:"foreignKey:UserID"`
	User       User       `json:"user "gorm:"foreignKey:UserID"`
	Credential Credential `json:"credential" gorm:"column:credential;type:json;not null"`
}

type Credential struct {
	ID              string                            `json:"id"`
	PublicKey       string                            `json:"publicKey"`
	AttestationType string                            `json:"attestationType"`
	Transport       []protocol.AuthenticatorTransport `json:"transport"`
	Flags           webauthn.CredentialFlags          `json:"flags"`
	Authenticator   webauthn.Authenticator            `json:"authenticator"`
}

func (c *Credential) Scan(value any) error {
	var cred Credential
	*c = cred
	if value == nil {
		return nil
	}

	var bytes []byte

	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	if err := json.Unmarshal(bytes, &cred); err != nil {
		return err
	}
	*c = cred
	return nil
}

func (c Credential) Value() (driver.Value, error) {
	bytes, err := json.Marshal(c)
	return string(bytes), err
}
