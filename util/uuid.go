package util

import (
	"strings"

	"github.com/google/uuid"
)

// GenerateUUID generates a pseudo-random UUID (version 4).
func GenerateUUID() string {
	s := uuid.New().String()
	return strings.ReplaceAll(s, "-", "")
}
