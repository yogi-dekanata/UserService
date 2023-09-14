package commons

import (
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	// EncryptionCost ...
	EncryptionCost = 16
	// SecretKey ...
	SecretKey = "&&hD23"
	// MaxPasswordLength ...
	MaxPasswordLength = 71
)

// PasswordManagerInterface this is contract
type PasswordManagerInterface interface {
	GenerateHash(password string, salt string) (string, error)
	VerifyPassword(password string, hash string, salt string) bool
	CreateSalt() string
}

// PasswordManager ...
type PasswordManager struct {
}

// GenerateHash ...
func (pm *PasswordManager) GenerateHash(rawPassword string, userSalt string) (string, error) {
	enhancedPassword := addSaltAndKey(rawPassword, userSalt)
	if len(enhancedPassword) > MaxPasswordLength {
		enhancedPassword = enhancedPassword[:MaxPasswordLength]
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(enhancedPassword), EncryptionCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword ...
func (pm *PasswordManager) VerifyPassword(rawPassword string, storedHash string, userSalt string) bool {
	userSalt = strings.TrimSpace(userSalt)

	enhancedPassword := addSaltAndKey(rawPassword, userSalt)
	if len(enhancedPassword) > MaxPasswordLength {
		enhancedPassword = enhancedPassword[:MaxPasswordLength]
	}

	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(enhancedPassword))
	if err != nil {
		return false
	}
	return true
}

// CreateSalt ...
func (pm *PasswordManager) CreateSalt() string {
	return strings.ReplaceAll(generateUUID(), "-", "")
}

// addSaltAndKey ...
func addSaltAndKey(password string, salt string) string {
	return SecretKey + password + salt
}

// generateUUID ...
func generateUUID() string {
	return uuid.New().String()
}
