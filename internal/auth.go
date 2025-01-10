package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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


func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration,) (string, error){
	claims := jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject: userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error){
	token, err := jwt.ParseWithClaims(tokenString,  &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error){
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid{
		return uuid.Nil, errors.New("invalid token")
	}
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}

	return  userID, nil
}

func GetBearerToken(headers http.Header) (string, error){
		authHeader := headers.Get("Authorization")
	if authHeader == ""{
		return "", errors.New("authorization header missing")
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer"{
		return "", errors.New("invalid authorzation header format")
	}

	return parts[1], nil

}
