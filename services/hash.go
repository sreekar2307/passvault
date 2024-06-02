package services

import (
	"context"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

type HashServiceImpl struct {
}

func NewHashService() HashServiceImpl {
	return HashServiceImpl{}
}

func (h HashServiceImpl) Hash(ctx context.Context, s string) (string, error) {
	passWordHash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(passWordHash), nil
}

func (h HashServiceImpl) CompareHash(ctx context.Context, stored string, userInput string) (bool, error) {
	hashedPassword, err := hex.DecodeString(stored)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(userInput))
	if err != nil {
		return false, err
	}
	return true, nil
}
