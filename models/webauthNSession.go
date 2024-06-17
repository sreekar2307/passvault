package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/google/uuid"
	"time"
)

type WebauthNSession struct {
	CreatedAt   time.Time   `json:"createdAt" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt   time.Time   `json:"updatedAt" gorm:"column:updated_at;type:timestamp;not null"`
	SessionID   uuid.UUID   `gorm:"column:session_id;default:uuid_generate_v4();type:uuid;not null;primaryKey" json:"sessionID"`
	UserID      uint        `json:"user_id "gorm:"foreignKey:UserID"`
	User        User        `json:"user "gorm:"foreignKey:UserID"`
	SessionData SessionData `json:"session" gorm:"column:session;type:json;not null"`
}

type SessionData struct {
	Challenge            string     `json:"challenge,omitempty"`
	UserID               string     `json:"user_id,omitempty"`
	AllowedCredentialIDs []string   `json:"allowed_credentials,omitempty"`
	Expires              *time.Time `json:"expires,omitempty"`

	UserVerification protocol.UserVerificationRequirement `json:"userVerification,omitempty"`
	Extensions       protocol.AuthenticationExtensions    `json:"extensions,omitempty"`
}

func (s *SessionData) Scan(value any) error {
	var sessionData SessionData
	*s = sessionData
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

	if err := json.Unmarshal(bytes, &sessionData); err != nil {
		return err
	}
	*s = sessionData
	return nil
}

func (s SessionData) Value() (driver.Value, error) {
	bytes, err := json.Marshal(s)
	return string(bytes), err
}
