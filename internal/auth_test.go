package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHasshPassword(t *testing.T) {
	password := "testPassword"

	hassh, err := HashPassword(password)
	if err != nil {
		fmt.Printf("Error :%s", err)
	}
	fmt.Printf("Hasshed Password: %s\n", hassh)
}

func TestPasswordLong(t *testing.T) {
	password := "testPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPassword"
	_, err := HashPassword(password)
	if err == nil { // If the hasshing succedding err would be nil, so failing the test
		t.Errorf("Should not have worked")
	}
}

func TestCheckPasswordHassh(t *testing.T) {
	hasshedPassword := "$2a$10$8u.iZ8zIGCJe1j3NIKI0eua3iS0JUV06qVP9pj/34p14dn36r1ZaK"
	password := "testPassword"

	err := CheckPasswordHash(password, hasshedPassword)
	if err != nil {
		fmt.Printf("Error :%s\n", err)
	}
}
func TestMakeAndValidateJWT(t *testing.T) {
	tokenSecret := "supersecretkey"
	userID := uuid.New()
	expiresIn := 1 * time.Hour
	// Generate the JWT
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate the JWT
	parsedUserID, err := ValidateJWT(token, tokenSecret)
	assert.NoError(t, err)
	assert.Equal(t, userID, parsedUserID)
}

func TestExpiredJWT(t *testing.T) {
	tokenSecret := "supersecretkey"
	userID := uuid.New()
	expiresIn := -1 * time.Second // Token already expired

	// Generate the JWT
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	assert.NoError(t, err)

	// Validate the JWT
	_, err = ValidateJWT(token, tokenSecret)
	assert.Error(t, err)
}

func TestInvalidSecret(t *testing.T) {
	tokenSecret := "supersecretkey"
	invalidSecret := "wrongsecret"
	userID := uuid.New()
	expiresIn := 1 * time.Hour

	// Generate the JWT
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	assert.NoError(t, err)

	// Validate the JWT with an invalid secret
	_, err = ValidateJWT(token, invalidSecret)
	assert.Error(t, err)
}
