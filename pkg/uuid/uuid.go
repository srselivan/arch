package uuid

import (
	"github.com/gofrs/uuid/v5"
)

func NewV7() string {
	return uuid.Must(uuid.NewV7()).String()
}
