package model

import "github.com/google/uuid"

type Secret struct {
	Data string
	ID   uuid.UUID
}
