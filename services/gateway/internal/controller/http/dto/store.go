package dto

import "github.com/google/uuid"

type CreateStoreRequest struct {
	Name    string    `json:"name"`
	OwnerID uuid.UUID `json:"owner_id,omitempty"`
}
