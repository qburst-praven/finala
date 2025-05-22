package serverutil

import (
	"crypto/rand"

	"math/big"
)

// GenerateRandomPassword creates a cryptographically secure random alphanumeric password.

// It uses crypto/rand for randomness, ensuring the password is hard to guess.

func GenerateRandomPassword(length int) (string, error) {

	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	if length <= 0 {

		length = 20 // Default length if invalid is provided

	}

	ret := make([]byte, length)

	for i := 0; i < length; i++ {

		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))

		if err != nil {

			return "", err // Propagate the error for the caller to handle

		}

		ret[i] = chars[num.Int64()]

	}

	return string(ret), nil

}
