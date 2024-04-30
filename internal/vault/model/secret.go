package model

import "github.com/google/uuid"

type Secret struct {
	Data []byte
	IV   []byte
	ID   uuid.UUID
}

func (s *Secret) Copy() *Secret {
	return &Secret{
		ID:   s.ID,
		IV:   s.IV,
		Data: s.Data,
	}
}
