package http

type Secret struct {
	ID   string `json:"id"`
	Data string `json:"data,omitempty"`
}

type ListSecretsResponse struct {
	List []Secret `json:"list,omitempty"`
}

type AddSecretRequest struct {
	Secret Secret `json:"secret"`
}

type AddSecretResponse struct {
	Secret Secret `json:"secret"`
}

type GetSecretResponse struct {
	Secret Secret `json:"secret"`
}
