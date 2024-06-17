package dtos

import (
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/google/uuid"
)

type GetWebAuthnCredentialFilter struct {
	CredID uuid.UUID
	UserID uint
}

type GetWebAuthnSessionFilter struct {
	SessionID uuid.UUID
	UserID    uint
}

type BeginRegisterResult struct {
	CredOptions *protocol.CredentialCreation `json:"credOptions"`
	SessionID   string                       `json:"sessionID"`
}

type BeginLoginResult struct {
	CredAssertion *protocol.CredentialAssertion `json:"credAssertion"`
	SessionID     string                        `json:"sessionID"`
}
