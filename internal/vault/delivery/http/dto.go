package http

type ListSecretsResponse struct {
	List []Secret `json:"list,omitempty"`
}

type Secret struct {
	ID      string `json:"id"`
	Payload string `json:"payload,omitempty"`
}

type AddSecretRequest struct {
	Secret Secret `json:"secret"`
}

type AddSecretResponse struct {
	Secret Secret `json:"secret"`
}
