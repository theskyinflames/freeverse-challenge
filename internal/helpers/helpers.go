package helpers

import (
	"github.com/google/uuid"
)

// BoolPtr is a helper
func BoolPtr(b bool) *bool {
	return &b
}

// UUIDPtr is a helper
func UUIDPtr(uuid uuid.UUID) *uuid.UUID {
	return &uuid
}
