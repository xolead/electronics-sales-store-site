package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Admin    bool
	Login    string
	RegData  time.Time
	Password string
}
