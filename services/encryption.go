package services

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"passVault/dtos"
	"passVault/interfaces"
	"passVault/models"
	"passVault/utils"
	"strings"
)

type EncryptionServiceImpl struct {
	encryption dtos.Encryption
	db         *gorm.DB
}

type SecretEncodingVersion string

var SecretEncodingVersions = struct {
	V1 SecretEncodingVersion
}{
	V1: "v1",
}
var SecretEncodingPattern = "%s|%s|%s"

func NewEncryptionService(db *gorm.DB, encryption dtos.Encryption) interfaces.EncryptService {
	return &EncryptionServiceImpl{db: db, encryption: encryption}
}

func (e EncryptionServiceImpl) Encrypt(_ context.Context, salt models.UserSalt, input string) (string, error) {
	encryptionKey := e.encryption.Keys[(utils.RandInt(len(e.encryption.Keys)))]
	ciphertext, err := e.encrypt(input, salt.Salt, encryptionKey.Key)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(SecretEncodingPattern, SecretEncodingVersions.V1, encryptionKey.Version, ciphertext), nil
}

func (e EncryptionServiceImpl) Decrypt(_ context.Context, salt models.UserSalt, input string) (string, error) {
	inputParts := strings.Split(input, "|")
	if len(inputParts) != 3 {
		return "", errors.New("invalid secret format")
	}

	version := inputParts[0]
	if version != string(SecretEncodingVersions.V1) {
		return "", errors.New("invalid secret encoding version")
	}
	keyVersion := inputParts[1]
	var encryptionKey string
	for _, key := range e.encryption.Keys {
		if key.Version == keyVersion {
			encryptionKey = key.Key
			break
		}
	}
	if encryptionKey == "" {
		return "", errors.New("invalid encryption key version")
	}

	ciphertext := inputParts[2]
	plainText, err := e.decrypt(ciphertext, salt.Salt, encryptionKey)
	if err != nil {
		return "", err
	}
	return string(plainText), nil
}

func (e EncryptionServiceImpl) EncryptBulk(ctx context.Context, salt models.UserSalt, inputs map[string]string) (map[string]string, error) {
	var encryptedInputs = make(map[string]string)
	for key, input := range inputs {
		encrypted, err := e.Encrypt(ctx, salt, input)
		if err != nil {
			return nil, err
		}
		encryptedInputs[key] = encrypted
	}
	return encryptedInputs, nil
}

func (e EncryptionServiceImpl) encrypt(plaintText, nonceAsHexString, keyAsHexString string) (string, error) {
	key, err := hex.DecodeString(keyAsHexString)
	if err != nil {
		return "", err
	}
	plainTextBytes := []byte(plaintText)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce, err := hex.DecodeString(nonceAsHexString)
	if err != nil {
		return "", err
	}
	nonce = nonce[:aesGCM.NonceSize()]
	ciphertext := aesGCM.Seal(nil, nonce, plainTextBytes, nil)
	return hex.EncodeToString(ciphertext), nil

}

func (e EncryptionServiceImpl) decrypt(cipherText, nonceAsHexString, keyAsHexString string) ([]byte, error) {
	key, err := hex.DecodeString(keyAsHexString)
	if err != nil {
		return nil, err
	}
	cipherTextBytes, err := hex.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce, err := hex.DecodeString(nonceAsHexString)
	if err != nil {
		return nil, err
	}
	nonce = nonce[:aesGCM.NonceSize()]
	plainText, err := aesGCM.Open(nil, nonce, cipherTextBytes, nil)
	return plainText, nil
}

func (e EncryptionServiceImpl) GenerateSalt(context.Context) (string, error) {
	var (
		randomBytes = make([]byte, 12)
	)
	_, err := utils.RandBytes(randomBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(randomBytes), nil
}
