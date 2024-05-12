package vault

import (
	httpVault "github.com/nestjam/goph-keeper/internal/vault/delivery/http"
)

type listSecretsCompletedMsg struct {
	secrets []*httpVault.Secret
}

type listSecretsFailedMsg struct {
	err        error
	statusCode int
}

type errMsg struct {
	err error
}

type getSecretCompletedMsg struct {
	secret httpVault.Secret
}

type getSecretFailedMsg struct {
	err        error
	secretID   string
	statusCode int
}

type deleteSecretCompletedMsg struct {
	secretID string
}

type deleteSecretFailedMsg struct {
	statusCode int
}

type createSecretRequestedMsg struct {
}

type saveSecretCompletedMsg struct {
	secret httpVault.Secret
}

type saveSecretFailedMsg struct {
	statusCode int
}
