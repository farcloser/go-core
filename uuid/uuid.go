package uuid

import (
	guuid "github.com/google/uuid"
)

func New() string {
	uuid, _ := guuid.NewV7()

	return uuid.String()
}
