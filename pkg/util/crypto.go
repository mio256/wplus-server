package util

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/taxio/errors"
)

func GeneratePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", errors.Wrap(err)
	}
	return string(hash), nil
}

func CompareHashAndPassword(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.Wrap(err)
	}
	return nil
}
