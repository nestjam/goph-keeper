package model

import "github.com/google/uuid"

type Secret struct {
	Payload string
	ID      uuid.UUID
}
