package helper

import "golang.org/x/crypto/bcrypt"

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

type BcryptHasher struct{}

func (b BcryptHasher) Compare(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (b BcryptHasher) Hash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}
