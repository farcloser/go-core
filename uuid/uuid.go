package uuid

import (
	guuid "github.com/google/uuid"
)

func New() string {
	guuid.EnableRandPool()
	uuid, _ := guuid.NewV7()
	return uuid.String()
}
