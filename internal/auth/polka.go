package auth

import (
	"errors"
	"net/http"
)

func GetAPIKey(headers http.Header) (string, error) {
	polkaKey := headers.Get("Authorization")
	if polkaKey == "" {
		return "", http.ErrNoCookie
	}

	if len(polkaKey) < 7 || polkaKey[:7] != "ApiKey " {
		return "", errors.New("invalid Authorization header format")
	}

	return polkaKey[7:], nil
}
