package vault

import "net/http"

type VaultHandlers interface {
	ListSecrets() http.HandlerFunc
	AddSecret() http.HandlerFunc
}
