package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const SessionTokenBytes = 32

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)

	nRead, err := rand.Read(b)

	if err != nil {
		return nil, fmt.Errorf("error reading random bytes: %w", err)
	}

	if nRead < n {
		return nil, fmt.Errorf("could not read enough random bytes: %w", err)
	}

	return b, nil
}

func String(n int) (string, error) {
	bytes, err := Bytes(n)

	if err != nil {
		return "", fmt.Errorf("error: %w", err)
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}

func SessionToken() (string, error) {
	return String(SessionTokenBytes)
}
