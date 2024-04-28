package http

type ListSecretsResponse struct {
	List []Secret `json:"list,omitempty"`
}

type Secret struct {
	ID   string `json:"id"`
	Data string `json:"data,omitempty"`
}

type AddSecretRequest struct {
	Secret Secret `json:"secret"`
}

type AddSecretResponse struct {
	Secret Secret `json:"secret"`
}
