package auth_service

import "errors"

var (
	ErrorInvalidCredentials = errors.New("invalid credentials")
	ErrorInvalidToken       = errors.New("invalid token")
)
