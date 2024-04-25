package config

import "time"

type JWTAuthConfig struct {
	SignKey       string
	TokenExpiryIn time.Duration
}
