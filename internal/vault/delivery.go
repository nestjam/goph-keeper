package vault

import "net/http"

type VaultHandlers interface {
	ListSecrets() http.HandlerFunc
	AddSecret() http.HandlerFunc
	UpdateSecret() http.HandlerFunc
	GetSecret() http.HandlerFunc
	DeleteSecret() http.HandlerFunc
}
