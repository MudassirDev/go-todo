package auth

import "golang.org/x/crypto/bcrypt"

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
