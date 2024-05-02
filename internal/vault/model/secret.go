package model

import "github.com/google/uuid"

type Secret struct {
	Data  []byte
	ID    uuid.UUID
	KeyID uuid.UUID
}

func (s *Secret) Copy() *Secret {
	return &Secret{
		ID:    s.ID,
		Data:  s.Data,
		KeyID: s.KeyID,
	}
}
