package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	ISSUER string = "todo"
)

func HashPassword(passwordString string) (string, error) {
	password := []byte(passwordString)

	rawHash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(rawHash)

	return hash, nil
}

func VerifyPassword(password, passwordHash string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}

func CreateJWT(jwtSecretString string, expiresIn time.Duration, userId string) (string, error) {
	jwtSecret := []byte(jwtSecretString)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    ISSUER,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   userId,
	})
	return token.SignedString(jwtSecret)
}

func VerifyJWT(jwtSecretString string, tokenString string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
		return []byte(jwtSecretString), nil
	})
	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}

	if issuer != ISSUER {
		return uuid.Nil, fmt.Errorf("invalid issuer")
	}

	rawId, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.Parse(rawId)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
