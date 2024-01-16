package hasher

import (
	"crypto/sha1"
	"fmt"
)

type passwordHasher struct {
	salt string
}

func NewPasswordHasher(salt string) Hasher {
	return &passwordHasher{
		salt: salt,
	}
}

func (h *passwordHasher) Hash(password string) (string, error) {
	hash := sha1.New()

	if _, err := hash.Write([]byte(password)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt))), nil
}
