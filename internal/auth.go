package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error){
	hasshed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hasshed), nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); 
	if err != nil {
		return err
	}
	return nil
}
