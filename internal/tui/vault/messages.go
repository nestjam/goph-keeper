package vault

import (
	httpVault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

type listSecretsCompletedMsg struct {
	secrets []httpVault.Secret
}

type listSecretsFailedMsg struct {
	statusCode int
}

type errMsg struct {
	err error
}

type getSecretCompletedMsg struct {
	secret httpVault.Secret
}

type getSecretFailedMsg struct {
	statusCode int
}

type deleteSecretCompletedMsg struct {
	secretID string
}

type deleteSecretFailedMsg struct {
	statusCode int
}
